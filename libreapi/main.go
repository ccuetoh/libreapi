package main

import "github.com/CamiloHernandez/libreapi/pkg"

func main() {
	certs := libreapi.TLSPaths{
		CertificatePath: "./certs/libreapi.pem",
		KeyPath:         "./certs/libreapi.key",
	}

	panic(libreapi.Start(443, certs))
}
