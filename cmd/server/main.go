package main

import (
	"github.com/ccuetoh/libreapi/pkg/config"
	libreapi "github.com/ccuetoh/libreapi/pkg/server"
	"log"

	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigFile("config.toml")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Unable to read config file: %v", err)
	}

	server, err := libreapi.NewServer(config.FromViper())
	if err != nil {
		log.Fatalf("Unable to start libreapi server: %v", err)
	}

	if err = server.Start(); err != nil {
		log.Fatalf("Unable to start libreapi server: %v", err)
	}
}
