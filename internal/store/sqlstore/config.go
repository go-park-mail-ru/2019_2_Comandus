package sqlstore

type Config struct {
	DatabaseURL string //`toml:"database_url"`
}

func NewConfig() *Config {
	return &Config{
		DatabaseURL:"host=localhost dbname=restapi_dev sslmode=disable port=5432 password=1234 user=d",
	}
}