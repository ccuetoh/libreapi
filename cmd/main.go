package main

import (
	"fmt"
	"os"

	"github.com/ccuetoh/libreapi/pkg/config"
	libreapi "github.com/ccuetoh/libreapi/pkg/router"

	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigFile("config.toml")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to read config file: %v", err)
		os.Exit(1)
	}

	server, err := libreapi.NewServer(config.FromViper())

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to start libreapi server: %v", err)
		os.Exit(1)
	}

	if err := server.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Unable to start libreapi server: %v", err)
		os.Exit(1)
	}
}
