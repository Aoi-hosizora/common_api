package main

import (
	"fmt"
	"github.com/Aoi-hosizora/common_api/api"
	"github.com/Aoi-hosizora/common_api/internal/pkg/module"
	"github.com/Aoi-hosizora/common_api/internal/server"
	"github.com/Aoi-hosizora/goapidoc"
	"github.com/spf13/pflag"
	"log"
)

var (
	fConfig = pflag.StringP("config", "c", "./config.yaml", "config file path")
	fHelp   = pflag.BoolP("help", "h", false, "show help")
)

func main() {
	pflag.Parse()
	if *fHelp {
		pflag.Usage()
		return
	}

	// module
	err := module.Provide(*fConfig)
	if err != nil {
		log.Fatalln("Failed to provide all modules:", err)
	}

	// document
	api.UpdateApiDoc()
	_, err = goapidoc.SaveSwaggerJson(api.SwaggerDocFilename)
	if err != nil {
		log.Fatalln("Failed to generate swagger:", err)
	}

	// server
	s, err := server.NewServer()
	if err != nil {
		log.Fatalln("Failed to create server:", err)
	}

	// start
	fmt.Println()
	s.Serve()
}
