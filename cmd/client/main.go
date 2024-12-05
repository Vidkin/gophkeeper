/*
Copyright © 2024 MIKHAIL SIRKIN <skim991@gmail.com>
*/

// Package main serves as the entry point for the GophKeeper application.
// It initializes and runs the command-line interface (CLI) using the cmd package,
// which handles all user commands and interactions. The package also provides
// build information such as version and date, which are printed at startup.
package main

import (
	"fmt"

	"github.com/Vidkin/gophkeeper/app/client/commands"
)

// buildVersion holds the build version of the application.
// It is set during the build process and defaults to "N/A" if not specified.
var buildVersion = "N/A"

// buildDate holds the build date of the application.
// It is set during the build process and defaults to "N/A" if not specified.
var buildDate = "N/A"

// main serves as the entry point for the GophKeeper application. It performs the following tasks:
//
//  1. Prints the build version and build date of the application.
//  2. Executes the command-line interface (CLI) using the commands.Execute() function,
//     which handles all user commands and interactions.
func main() {
	fmt.Printf("Build version: %s\nBuild date: %s\n---------------\n", buildVersion, buildDate)
	commands.Execute()
}
