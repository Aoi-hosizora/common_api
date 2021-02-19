package main

import (
	"flag"
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

	_, err := goapidoc.GenerateSwaggerJson("./api/spec.json")
	if err != nil {
		log.Fatalln("Failed to generate swagger:", err)
	}
	_, err = goapidoc.GenerateApib("./api/spec.apib")
	if err != nil {
		log.Fatalln("Failed to generate api blueprint:", err)
	}
	err = module.Provide(*fConfig)
	if err != nil {
		log.Fatalln("Failed to provide some modules:", err)
	}

	s, err := server.NewServer()
	if err != nil {
		log.Fatalln("Failed to create server:", err)
	}
	s.Serve()
}
