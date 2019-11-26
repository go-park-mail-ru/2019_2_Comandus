package apiserver

import (
	"database/sql"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/store/create"
	"github.com/go-park-mail-ru/2019_2_Comandus/monitoring"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
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
	prometheus.MustRegister(monitoring.FooCount, monitoring.Hits)

	config := NewConfig()

	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = ":8080"
	} else {
		port = ":" + port
	}
	config.BindAddr = port

	url :=  os.Getenv("DATABASE_URL")
	if len(url) != 0 {
		config.DatabaseURL = url
	}

	zapLogger, _ := zap.NewProduction()
	defer func() {
		if err := zapLogger.Sync(); err != nil {
			log.Println("HEHEHEG",err)
		}
	}()
	sugaredLogger := zapLogger.Sugar()

	srv, err := NewServer(config, sugaredLogger)
	if err != nil {
		return err
	}

	db, err := newDB(config.DatabaseURL)
	if err != nil {
		return err
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Println(err)
		}
	}()



	srv.ConfigureServer(db)

	/*if err := httpscerts.Check("cert.pem", "key.pem"); err != nil {
		err = httpscerts.Generate("cert.pem", "key.pem", "127.0.0.1:8081")
		if err != nil {
			log.Fatal("Ошибка: Не можем сгенерировать https сертификат.")
		}
	}
	return http.ListenAndServeTLS(":8081", "cert.pem", "key.pem", nil)*/
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
