package apiserver

import "github.com/go-park-mail-ru/2019_2_Comandus/internal/store"

type Config struct {
	BindAddr    string
	LogLevel    string
	DatabaseURL string
	SessionKey  string
	Store *store.Config
}

// NewConfig ...
func NewConfig() *Config {
	return &Config{
		BindAddr: ":8080",
		LogLevel: "debug",
		SessionKey: "jdfhdfdj",
		DatabaseURL: "host=localhost dbname=flapi sslmode=disable", // ??
		//Store: store.NewConfig(),
	}
}
