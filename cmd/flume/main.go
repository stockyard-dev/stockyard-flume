package main
import ("fmt";"log";"os";"github.com/stockyard-dev/stockyard-flume/internal/server";"github.com/stockyard-dev/stockyard-flume/internal/store")
func main(){port:=os.Getenv("PORT");if port==""{port="9210"};dataDir:=os.Getenv("DATA_DIR");if dataDir==""{dataDir=", "}
db,err:=store.Open(dataDir);if err!=nil{log.Fatalf("flume: %v",err)};defer db.Close();srv:=server.New(db,server.DefaultLimits())
fmt.Printf("\n  Stockyard Flume\n  Dashboard:  http://localhost:%s/ui\n  API:        http://localhost:%s/api\n  Questions? hello@stockyard.dev — I read every message\n\n",port,port)
log.Printf("flume: listening on :%s",port);log.Fatal(srv.ListenAndServe(":"+port))}
