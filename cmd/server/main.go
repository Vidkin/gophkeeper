package main

import (
	"fmt"

	"github.com/Vidkin/gophkeeper/app"
	"github.com/Vidkin/gophkeeper/internal/config"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
)

func main() {
	fmt.Printf("Build version: %s\nBuild date: %s\n", buildVersion, buildDate)

	cfg, err := config.NewServerConfig()
	if err != nil {
		panic(err)
	}

	serverApp, err := app.NewServerApp(cfg)
	if err != nil {
		panic(err)
	}
	fmt.Println(serverApp)
	// TODO: serverApp.Run()
}
