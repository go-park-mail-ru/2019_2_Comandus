package apiserver

type Config struct {
	BindAddr    string
	LogLevel    string
	//DatabaseURL string
	SessionKey  string
}

// NewConfig ...
func NewConfig() *Config {
	return &Config{
		BindAddr: ":8080",
		LogLevel: "debug",
		SessionKey: "jdfhdfdj",
	}
}
