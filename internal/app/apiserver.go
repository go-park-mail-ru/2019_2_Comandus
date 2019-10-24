package apiserver

import (
	"database/sql"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/store/sqlstore"
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
	db, err := newDB(config.DatabaseURL)
	if err != nil {
		return err
	}
	defer db.Close()

	store := sqlstore.New(db)

	sessionStore := sessions.NewCookieStore([]byte(config.SessionKey))
	srv := newServer(sessionStore, store)

	return http.ListenAndServe(config.BindAddr, srv)
}

func newDB(dbURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, err
	}

	// MAX CONNECTIONS
	db.SetMaxOpenConns(10)

	/*if err := db.Ping(); err != nil {
		return nil, err
	}*/
	if err := createTables(db); err != nil {
		return nil, err
	}
	return db, nil
}

func createTables(db *sql.DB) error {
	usersQuery := `CREATE TABLE IF NOT EXISTS users (
		accountId bigserial not null primary key,
		firstName varchar,
		secondName varchar,
		userName varchar not null,
		email varchar not null unique,
		encryptPassword varchar not null,
		avatar bytea,
		userType varchar --not null
	);`
	if _, err := db.Exec(usersQuery); err != nil {
		return err
	}


	managersQuery := `CREATE TABLE IF NOT EXISTS managers (
		id bigserial not null primary key,
		accountId bigserial references users,
		registrationDate timestamp,
		location varchar,
		companyId bigserial --references  companies
	);`
	if _, err := db.Exec(managersQuery); err != nil {
		return err
	}

	freelancersQuery := `CREATE TABLE IF NOT EXISTS freelancers (
		id bigserial not null primary key,
		accountId bigserial not null references users,
		registrationDate timestamp,
		country varchar,
		city varchar,
		address varchar,
		phone varchar,
		tagLine varchar,
		overview varchar,
		experienceLevelId bigserial,
		specialityId bigserial --references specialities
	);`
	if _, err := db.Exec(freelancersQuery); err != nil {
		return err
	}

	return nil
}