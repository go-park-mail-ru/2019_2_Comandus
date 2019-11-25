package apiserver

type Config struct {
	BindAddr    string
	LogLevel    string
	DatabaseURL string
	SessionKey  string
	TokenSecret string
	ClientUrl	string
}

func NewConfig() *Config {
	return &Config{
		BindAddr:		":8080",
		LogLevel:		"debug",
		SessionKey:		"jdfhdfdj",
		DatabaseURL:	"host=localhost dbname=restapi_dev sslmode=disable port=5432 password=1234 user=d",
		TokenSecret:	"golangsecpark",
		//ClientUrl:		"http://127.0.0.1:9000",
		ClientUrl:		"https://comandus.now.sh",
	}
}
