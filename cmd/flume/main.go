package main

import (
	"flag"
	"fmt"
	"github.com/stockyard-dev/stockyard-flume/internal/server"
	"github.com/stockyard-dev/stockyard-flume/internal/store"
	"log"
	"os"
)

func main() {
	portFlag := flag.String("port", "", "")
	dataFlag := flag.String("data", "", "")
	flag.Parse()
	port := os.Getenv("PORT")
	if port == "" {
		port = "9210"
	}
	dataDir := os.Getenv("DATA_DIR")
	if dataDir == "" {
		dataDir = ", "
	}
	if *portFlag != "" {
		port = *portFlag
	}
	if *dataFlag != "" {
		dataDir = *dataFlag
	}
	db, err := store.Open(dataDir)
	if err != nil {
		log.Fatalf("flume: %v", err)
	}
	defer db.Close()
	srv := server.New(db, server.DefaultLimits(), dataDir)
	fmt.Printf("\n  Stockyard Flume\n  Dashboard:  http://localhost:%s/ui\n  API:        http://localhost:%s/api\n  Questions? hello@stockyard.dev — I read every message\n\n", port, port)
	log.Printf("flume: listening on :%s", port)
	log.Fatal(srv.ListenAndServe(":" + port))
}
