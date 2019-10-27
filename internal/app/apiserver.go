package apiserver

import (
	"database/sql"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/store/sqlstore"
	"github.com/gorilla/csrf"
	"github.com/gorilla/sessions"
	"go.uber.org/zap"
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
	zapLogger, _ := zap.NewProduction()
	sugaredLogger := zapLogger.Sugar()
	db, err := newDB(config.DatabaseURL)
	if err != nil {
		return err
	}
	defer db.Close()
	defer zapLogger.Sync()

	store := sqlstore.New(db)

	sessionStore := sessions.NewCookieStore([]byte(config.SessionKey))
	srv := newServer(sessionStore, store, sugaredLogger)

	CSRF := csrf.Protect([]byte("32-byte-long-auth-key"))
	return http.ListenAndServe(config.BindAddr, CSRF(srv))
}

func newDB(dbURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, err
	}

	// MAX CONNECTIONS
	db.SetMaxOpenConns(20)

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

	jobsQuery := `CREATE TABLE IF NOT EXISTS jobs (
		id bigserial not null primary key,
		managerId bigserial not null references managers,
		title varchar not null,
		description varchar not null,
		files varchar,
		specialityId bigserial,  --references specialities,
		experienceLevelId bigserial,
		paymentAmount float8,
		country varchar,
		city varchar,
		jobTypeId bigserial
	);`
	if _, err := db.Exec(jobsQuery); err != nil {
		return err
	}

	specialitiesQuery := `CREATE TABLE IF NOT EXISTS specialities (
		id bigserial not null primary key,
		name varchar
	);`
	if _, err := db.Exec(specialitiesQuery); err != nil {
		return err
	}

	companiesQuery := `CREATE TABLE IF NOT EXISTS companies (
		id bigserial not null primary key,
		name varchar
	);`

	if _, err := db.Exec(companiesQuery); err != nil {
		return err
	}

	return nil
}

func dropAllTables(db *sql.DB) error {
	query := `DROP TABLE IF EXISTS users;`
	if _, err := db.Exec(query); err != nil {
		return err
	}

	query = `DROP TABLE IF EXISTS managers;`
	if _, err := db.Exec(query); err != nil {
		return err
	}

	query = `DROP TABLE IF EXISTS freelancers;`
	if _, err := db.Exec(query); err != nil {
		return err
	}

	query = `DROP TABLE IF EXISTS jobs;`
	if _, err := db.Exec(query); err != nil {
		return err
	}

	query = `DROP TABLE IF EXISTS companies;`
	if _, err := db.Exec(query); err != nil {
		return err
	}

	query = `DROP TABLE IF EXISTS specialities;`
	if _, err := db.Exec(query); err != nil {
		return err
	}

	return nil
}
