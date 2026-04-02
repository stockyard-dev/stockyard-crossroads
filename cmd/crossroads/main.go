package main
import ("fmt";"log";"net/http";"os";"github.com/stockyard-dev/stockyard-crossroads/internal/server";"github.com/stockyard-dev/stockyard-crossroads/internal/store")
func main(){port:=os.Getenv("PORT");if port==""{port="9070"};dataDir:=os.Getenv("DATA_DIR");if dataDir==""{dataDir="./crossroads-data"}
db,err:=store.Open(dataDir);if err!=nil{log.Fatalf("crossroads: %v",err)};defer db.Close();srv:=server.New(db)
fmt.Printf("\n  Crossroads — Self-hosted URL shortener\n  ─────────────────────────────────\n  Dashboard:  http://localhost:%s/ui\n  API:        http://localhost:%s/api\n  Data:       %s\n  ─────────────────────────────────\n\n",port,port,dataDir)
log.Printf("crossroads: listening on :%s",port);log.Fatal(http.ListenAndServe(":"+port,srv))}
