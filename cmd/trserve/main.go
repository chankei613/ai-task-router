// cmd/trserve はAI Task Router APIをlocalhostで提供する単体サーバー。
//
//	go run ./cmd/trserve -addr :8427 -db ai-task-router.db
package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/chankei613/ai-task-router/internal/api"
	"github.com/chankei613/ai-task-router/internal/db"
)

func main() {
	addr := flag.String("addr", ":8427", "待ち受けアドレス")
	dbPath := flag.String("db", "ai-task-router.db", "SQLiteファイル")
	flag.Parse()

	conn, err := db.Init(*dbPath)
	if err != nil {
		log.Fatalf("db init failed: %v", err)
	}

	router := api.NewRouter(conn)
	log.Printf("ai-task-router backend listening on %s", *addr)
	if err := http.ListenAndServe(*addr, router); err != nil {
		log.Fatal(err)
	}
}
