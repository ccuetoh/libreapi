package main

import (
	"fmt"
	"github.com/ccuetoh/libreapi/pkg/config"
	libreapi "github.com/ccuetoh/libreapi/pkg/router"
	"os"
)

func main() {
	server, err := libreapi.NewServer(
		config.SetNewRelicLicence(""),
		config.SetPort(8080))

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to start libreapi server: %v", err)
		os.Exit(1)
	}

	if err := server.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Unable to start libreapi server: %v", err)
		os.Exit(1)
	}
}
