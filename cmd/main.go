package main

import (
	"flag"
	"fmt"
	"github.com/Aoi-hosizora/common_api/api"
	"github.com/Aoi-hosizora/common_api/internal/pkg/module"
	"github.com/Aoi-hosizora/common_api/internal/server"
	"github.com/Aoi-hosizora/goapidoc"
	"log"
)

var (
	fConfig = flag.String("config", "./config.yaml", "config file path")
	fHelp   = flag.Bool("h", false, "show help")
)

func main() {
	flag.Parse()
	if *fHelp {
		flag.Usage()
		return
	}

	err := module.Provide(*fConfig)
	if err != nil {
		log.Fatalln("Failed to provide some of the modules:", err)
	}

	api.UpdateApiDoc()
	_, err = goapidoc.SaveSwaggerJson(api.SwaggerDocFilename)
	if err != nil {
		log.Fatalln("Failed to generate swagger:", err)
	}

	s, err := server.NewServer()
	if err != nil {
		log.Fatalln("Failed to create server:", err)
	}
	fmt.Println()
	s.Serve()
}
