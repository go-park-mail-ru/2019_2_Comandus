package main

import (
	"flag"
	apiserver "github.com/go-park-mail-ru/2019_2_Comandus/internal/app"
	"log"
	"os"
)

var (
	configPath         string
	localhostClientUrl string = "http://localhost:9000"
	runForLocalClient  bool
)

func init() {
	flag.StringVar(&configPath, "config-path", "configs/apiserver.toml", "path to config file")
	flag.BoolVar(&runForLocalClient, "local", false, "assign clientUrl as "+localhostClientUrl)
}

func main() {
	flag.Parse()
	config := apiserver.NewConfig()
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = ":8080"
	} else {
		port = ":" + port
	}
	config.BindAddr = port

	url := os.Getenv("DATABASE_URL")
	if len(url) != 0 {
		config.DatabaseURL = url
	}

	if runForLocalClient {
		config.ClientUrl = localhostClientUrl
	}

	if err := apiserver.Start(config); err != nil {
		log.Fatal(err)
	}
}
