package apiserver

import (
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/store/sqlstore"
)

type Config struct {
	BindAddr    string
	LogLevel    string
	DatabaseURL string
	SessionKey  string
	Store       *sqlstore.Config
	TokenSecret string
}

func NewConfig() *Config {
	return &Config{
		BindAddr:    ":8080",
		LogLevel:    "debug",
		SessionKey:  "jdfhdfdj",
		DatabaseURL: "host=localhost dbname=restapi_dev sslmode=disable port=5432 password=1234 user=d",
		Store:       sqlstore.NewConfig(),
		TokenSecret: "golangsecpark",
	}
}
