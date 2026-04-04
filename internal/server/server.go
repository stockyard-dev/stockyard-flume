package server
import("encoding/json";"net/http";"github.com/stockyard-dev/stockyard-flume/internal/store")
type Server struct{db *store.DB;limits Limits;mux *http.ServeMux}
func New(db *store.DB,tier string)*Server{s:=&Server{db:db,limits:LimitsFor(tier),mux:http.NewServeMux()};s.routes();return s}
func(s *Server)ListenAndServe(addr string)error{return(&http.Server{Addr:addr,Handler:s.mux}).ListenAndServe()}
func(s *Server)routes(){
    s.mux.HandleFunc("GET /health",s.handleHealth)
    s.mux.HandleFunc("GET /api/stats",s.handleStats)
    s.mux.HandleFunc("GET /api/streams",s.handleListStreams)
    s.mux.HandleFunc("POST /api/streams",s.handleCreateStream)
    s.mux.HandleFunc("DELETE /api/streams/{id}",s.handleDeleteStream)
    s.mux.HandleFunc("GET /api/logs",s.handleQueryLogs)
    s.mux.HandleFunc("POST /api/ingest/{id}",s.handleIngest)
    s.mux.HandleFunc("GET /",s.handleUI)
s.mux.HandleFunc("GET /api/tier",func(w http.ResponseWriter,r *http.Request){writeJSON(w,200,map[string]any{"tier":s.limits.Tier,"upgrade_url":"https://stockyard.dev/flume/"})})
}
func(s *Server)handleHealth(w http.ResponseWriter,r *http.Request){writeJSON(w,200,map[string]string{"status":"ok","service":"stockyard-flume"})}
func writeJSON(w http.ResponseWriter,status int,v interface{}){w.Header().Set("Content-Type","application/json");w.WriteHeader(status);json.NewEncoder(w).Encode(v)}
func writeError(w http.ResponseWriter,status int,msg string){writeJSON(w,status,map[string]string{"error":msg})}
func(s *Server)handleUI(w http.ResponseWriter,r *http.Request){if r.URL.Path!="/"{http.NotFound(w,r);return};w.Header().Set("Content-Type","text/html");w.Write(dashboardHTML)}
