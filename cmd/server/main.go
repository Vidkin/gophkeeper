package main

import (
	"fmt"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
)

func main() {
	fmt.Printf("Build version: %s\nBuild date: %s\n", buildVersion, buildDate)

	//cfg, err := config.NewServerConfig()
	//if err != nil {
	//	panic(err)
	//}
	//
	//serverApp, err := app.NewServerApp(cfg)
	//if err != nil {
	//	panic(err)
	//}
	//serverApp.Run()
}
