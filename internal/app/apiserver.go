package apiserver

import (
	"database/sql"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/store/create"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/microcosm-cc/bluemonday"
	"go.uber.org/zap"
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

func Start(config *Config) error {
	zapLogger, _ := zap.NewProduction()
	defer func() {
		if err := zapLogger.Sync(); err != nil {
			log.Println(err)
		}
	}()

	sugaredLogger := zapLogger.Sugar()

	db, err := newDB(config.DatabaseURL)
	if err != nil {
		return err
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Println(err)
		}
	}()

	token, err := NewHMACHashToken(config.TokenSecret)
	if err != nil {
		return err
	}

	sanitizer := bluemonday.UGCPolicy()
	sessionStore := sessions.NewCookieStore([]byte(config.SessionKey))

	m := mux.NewRouter()
	srv := NewServer(m, sessionStore, sugaredLogger, token, sanitizer)
	srv.ConfigureServer(db)
	return http.ListenAndServe(config.BindAddr, srv)
}

func newDB(dbURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(20)
	if err := create.CreateTables(db); err != nil {
		return nil, err
	}
	return db, nil
}