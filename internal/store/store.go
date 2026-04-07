package store

import (
	"database/sql"
	"fmt"
	_ "modernc.org/sqlite"
	"os"
	"path/filepath"
	"time"
)

type DB struct{ *sql.DB }
type Stream struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Token     string    `json:"token"`
	Retention int       `json:"retention_days"`
	CreatedAt time.Time `json:"created_at"`
}
type LogEntry struct {
	ID         int64     `json:"id"`
	StreamID   int64     `json:"stream_id"`
	StreamName string    `json:"stream_name,omitempty"`
	Level      string    `json:"level"`
	Message    string    `json:"message"`
	Fields     string    `json:"fields"`
	CreatedAt  time.Time `json:"created_at"`
}

func Open(dataDir string) (*DB, error) {
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, fmt.Errorf("mkdir: %w", err)
	}
	dsn := filepath.Join(dataDir, "flume.db") + "?_journal_mode=WAL&_busy_timeout=5000"
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("open: %w", err)
	}
	db.SetMaxOpenConns(1)
	if err := migrate(db); err != nil {
		return nil, fmt.Errorf("migrate: %w", err)
	}
	db.Exec(`CREATE TABLE IF NOT EXISTS extras(resource TEXT NOT NULL,record_id TEXT NOT NULL,data TEXT NOT NULL DEFAULT '{}',PRIMARY KEY(resource, record_id))`)
	return &DB{db}, nil
}
func migrate(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS streams(id INTEGER PRIMARY KEY AUTOINCREMENT,name TEXT NOT NULL UNIQUE,token TEXT NOT NULL,retention_days INTEGER DEFAULT 7,created_at DATETIME DEFAULT CURRENT_TIMESTAMP);CREATE TABLE IF NOT EXISTS log_entries(id INTEGER PRIMARY KEY AUTOINCREMENT,stream_id INTEGER NOT NULL,level TEXT DEFAULT 'info',message TEXT NOT NULL,fields TEXT DEFAULT '{}',created_at DATETIME DEFAULT CURRENT_TIMESTAMP);CREATE INDEX IF NOT EXISTS log_stream ON log_entries(stream_id,created_at DESC);CREATE INDEX IF NOT EXISTS log_level ON log_entries(level);`)
	return err
}
func (db *DB) ListStreams() ([]Stream, error) {
	rows, err := db.Query(`SELECT id,name,token,retention_days,created_at FROM streams ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Stream
	for rows.Next() {
		var s Stream
		rows.Scan(&s.ID, &s.Name, &s.Token, &s.Retention, &s.CreatedAt)
		out = append(out, s)
	}
	return out, nil
}
func (db *DB) CreateStream(s *Stream) error {
	if s.Token == "" {
		s.Token = fmt.Sprintf("%x", time.Now().UnixNano())
	}
	if s.Retention == 0 {
		s.Retention = 7
	}
	res, err := db.Exec(`INSERT INTO streams(name,token,retention_days)VALUES(?,?,?)`, s.Name, s.Token, s.Retention)
	if err != nil {
		return err
	}
	s.ID, _ = res.LastInsertId()
	return nil
}
func (db *DB) DeleteStream(id int64) error {
	_, err := db.Exec(`DELETE FROM streams WHERE id=?`, id)
	_, _ = db.Exec(`DELETE FROM log_entries WHERE stream_id=?`, id)
	return err
}
func (db *DB) QueryLogs(streamID int64, level, search string, limit int) ([]LogEntry, error) {
	if limit == 0 || limit > 1000 {
		limit = 200
	}
	q := `SELECT l.id,l.stream_id,COALESCE(s.name,''),l.level,l.message,l.fields,l.created_at FROM log_entries l LEFT JOIN streams s ON s.id=l.stream_id WHERE 1=1`
	var args []interface{}
	if streamID > 0 {
		q += " AND l.stream_id=?"
		args = append(args, streamID)
	}
	if level != "" {
		q += " AND l.level=?"
		args = append(args, level)
	}
	if search != "" {
		q += " AND l.message LIKE ?"
		args = append(args, "%"+search+"%")
	}
	q += fmt.Sprintf(" ORDER BY l.created_at DESC LIMIT %d", limit)
	rows, err := db.Query(q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []LogEntry
	for rows.Next() {
		var e LogEntry
		rows.Scan(&e.ID, &e.StreamID, &e.StreamName, &e.Level, &e.Message, &e.Fields, &e.CreatedAt)
		out = append(out, e)
	}
	return out, nil
}
func (db *DB) IngestLog(streamID int64, level, message, fields string) error {
	if level == "" {
		level = "info"
	}
	if fields == "" {
		fields = "{}"
	}
	_, err := db.Exec(`INSERT INTO log_entries(stream_id,level,message,fields)VALUES(?,?,?,?)`, streamID, level, message, fields)
	return err
}
func (db *DB) PurgeLogs(streamID, retentionDays int64) {
	db.Exec(`DELETE FROM log_entries WHERE stream_id=? AND created_at < datetime('now','-'||?||' days')`, streamID, retentionDays)
}
func (db *DB) CountLogs(streamID int64) (int, error) {
	var n int
	if streamID > 0 {
		db.QueryRow(`SELECT COUNT(*) FROM log_entries WHERE stream_id=?`, streamID).Scan(&n)
	} else {
		db.QueryRow(`SELECT COUNT(*) FROM log_entries`).Scan(&n)
	}
	return n, nil
}

// ─── Extras: generic key-value storage for personalization custom fields ───

func (d *DB) GetExtras(resource, recordID string) string {
	var data string
	err := d.QueryRow(
		`SELECT data FROM extras WHERE resource=? AND record_id=?`,
		resource, recordID,
	).Scan(&data)
	if err != nil || data == "" {
		return "{}"
	}
	return data
}

func (d *DB) SetExtras(resource, recordID, data string) error {
	if data == "" {
		data = "{}"
	}
	_, err := d.Exec(
		`INSERT INTO extras(resource, record_id, data) VALUES(?, ?, ?)
		 ON CONFLICT(resource, record_id) DO UPDATE SET data=excluded.data`,
		resource, recordID, data,
	)
	return err
}

func (d *DB) DeleteExtras(resource, recordID string) error {
	_, err := d.Exec(
		`DELETE FROM extras WHERE resource=? AND record_id=?`,
		resource, recordID,
	)
	return err
}

func (d *DB) AllExtras(resource string) map[string]string {
	out := make(map[string]string)
	rows, _ := d.Query(
		`SELECT record_id, data FROM extras WHERE resource=?`,
		resource,
	)
	if rows == nil {
		return out
	}
	defer rows.Close()
	for rows.Next() {
		var id, data string
		rows.Scan(&id, &data)
		out[id] = data
	}
	return out
}
