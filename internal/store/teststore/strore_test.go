package teststore

import (
	"database/sql"
	"fmt"
	"strings"
	"testing"
)

const (
	databaseURL = "host=localhost dbname=restapi_dev sslmode=disable port=5432 password=1234 user=d"
)

func testStore(t *testing.T, databaseURL string) (*sql.DB, func(...string)) {
	t.Helper()
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		t.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		t.Fatal(err)
	}

	return db, func(tables ...string) {
		if len(tables) > 0 {
			_, err = db.Exec(fmt.Sprintf("TRUNCATE %s CASCADE", strings.Join(tables, ", ")))
			if err != nil {
				t.Fatal(err)
			}
		}
		if err := db.Close(); err != nil {
			t.Fatal(err)
		}
	}
}
