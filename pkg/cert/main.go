/*
Package main provides the entry point for creating and saving
X.509 certificates and corresponding private keys in PEM format.
*/
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Vidkin/gophkeeper/pkg/cert/x509"
)

func main() {
	if len(os.Args) != 5 {
		fmt.Println("Usage: go run main.go <organization> <country> <pathCertPEM> <pathPrivateKeyPEM>")
		return
	}

	organization := os.Args[1]
	country := os.Args[2]
	pathCertPEM := os.Args[3]
	pathPrivateKeyPEM := os.Args[4]

	err := x509.CreateAndSave(organization, country, pathCertPEM, pathPrivateKeyPEM)
	if err != nil {
		log.Fatalf("Error creating and saving certificate: %v", err)
	}

	fmt.Println("Certificate and private key created and saved successfully.")
}
