/*
Copyright © 2024 MIKHAIL SIRKIN <skim991@gmail.com>
*/
package main

import (
	"fmt"

	"github.com/Vidkin/gophkeeper/app/cmd"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
)

func main() {
	fmt.Printf("Build version: %s\nBuild date: %s\n---------------\n", buildVersion, buildDate)
	cmd.Execute()

}
