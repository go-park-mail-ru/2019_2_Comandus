package main

import (
	"flag"
	apiserver "github.com/go-park-mail-ru/2019_2_Comandus/internal/app"
	"log"
)

var (
	configPath         string
	localhostClientUrl string = "http://localhost:9000"
	runForLocalClient  bool
)

func init() {
	flag.StringVar(&configPath, "config-path", "configs/apiserver.toml", "path to config file")
	flag.BoolVar(&runForLocalClient, "local", false, "assign clientUrl as " + localhostClientUrl)
}

func main() {
	//flag.Parse()
	if err := apiserver.Start(); err != nil {
		log.Fatal(err)
	}
}
