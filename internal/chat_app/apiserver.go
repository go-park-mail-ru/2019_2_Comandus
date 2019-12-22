package chat_app

import (
	"database/sql"
	"flag"
	"fmt"
	store "github.com/go-park-mail-ru/2019_2_Comandus/internal/store/create"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

type responseWriter struct {
	http.ResponseWriter
	code int
}

func (w *responseWriter) WriteHeader(statusCode int) {
	w.code = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func Start() error {
	flag.Parse()

	config := NewConfig()
	db, _ := newDB(config.DatabaseURL)

	log.SetFlags(log.Lshortfile)

	// websocket server
	server := NewServer("/entry", db)
	go server.Listen()

	// static files
	http.Handle("/", http.FileServer(http.Dir("webroot")))

	fmt.Println("starting chat_app at :8089")
	return http.ListenAndServe(":8089", nil)
}

func newDB(dbURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(20)
	if err := store.CreateChatTables(db); err != nil {
		return nil, err
	}
	return db, nil
}
