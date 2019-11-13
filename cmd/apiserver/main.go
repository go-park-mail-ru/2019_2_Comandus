package main

import (
	"flag"
	apiserver "github.com/go-park-mail-ru/2019_2_Comandus/internal/app"
	"log"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "configs/apiserver.toml", "path to config file")
}

func main() {
	//flag.Parse()
	if err := apiserver.Start(); err != nil {
		log.Fatal(err)
	}
}