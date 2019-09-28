package apiserver

import (
	"database/sql"
	"github.com/gorilla/sessions"
	"net/http"
)

// Зачем обертки над http.ResponseWriter , код можно и в нем засетить
type responseWriter struct {
	http.ResponseWriter
	code int
}

func (w *responseWriter) WriteHeader(statusCode int) {
	w.code = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func Start(config *Config) error {
	/*db, err := newDB(config.DatabaseURL)
	if err != nil {
		return err
	}

	defer db.Close()
	store := store.New(db)*/


	sessionStore := sessions.NewCookieStore([]byte(config.SessionKey))
	srv := newServer(sessionStore)
	//err := srv.ConfigureStore()
	//if err != nil {
	//	return err
	//}
	return http.ListenAndServe(config.BindAddr, srv)
}
func newDB(dbURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}