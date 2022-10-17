package main

import (
	"fmt"
	"os"

	libreapi "github.com/ccuetoh/libreapi/pkg"
)

func main() {
	certs := libreapi.TLSPaths{
		CertificatePath: "./certs/libreapi.pem",
		KeyPath:         "./certs/libreapi.key",
	}

	if err := libreapi.Start(443, certs); err != nil {
		fmt.Fprintf(os.Stderr, "Unable to start libreapi server: %v", err)
		os.Exit(1)
	}
}
