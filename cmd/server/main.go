/*
Copyright © 2024 MIKHAIL SIRKIN <skim991@gmail.com>
*/

// Package main serves as the entry point for the GophKeeper application.
// It initializes and runs the server using configurations provided by the config package,
// and handles server operations through the app package.
package main

import (
	"fmt"

	"github.com/Vidkin/gophkeeper/app"
	"github.com/Vidkin/gophkeeper/internal/config"
)

// buildVersion holds the build version of the application.
// It is set during the build process and defaults to "N/A" if not specified.
var buildVersion = "N/A"

// buildDate holds the build date of the application.
// It is set during the build process and defaults to "N/A" if not specified.
var buildDate = "N/A"

// main serves as the entry point for the GophKeeper server. It performs the following tasks:
//
// 1. Prints the build version and build date of the application.
// 2. Loads the server configuration using config.NewServerConfig().
// 3. Initializes the server application using app.NewServerApp(cfg).
// 4. Starts the server by calling serverApp.Run().
func main() {
	fmt.Printf("Build version: %s\nBuild date: %s\n---------------\n", buildVersion, buildDate)
	cfg, err := config.NewServerConfig()
	if err != nil {
		panic(err)
	}
	serverApp, err := app.NewServerApp(cfg)
	if err != nil {
		panic(err)
	}
	serverApp.Run()
}
