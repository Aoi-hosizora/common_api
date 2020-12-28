package main

import (
	"flag"
	"github.com/Aoi-hosizora/common_api/src/provide"
	"github.com/Aoi-hosizora/common_api/src/server"
	"github.com/Aoi-hosizora/goapidoc"
	"log"
)

var (
	fHelp   = flag.Bool("h", false, "show help")
	fConfig = flag.String("config", "./config.yaml", "set config file")
)

func main() {
	if *fHelp {
		flag.Usage()
	} else {
		run()
	}
}

func run() {
	_, err := goapidoc.GenerateSwaggerJson("./docs/doc.json")
	if err != nil {
		log.Fatalln("Failed to generate swagger:", err)
	}
	_, err = goapidoc.GenerateApib("./docs/doc.apib")
	if err != nil {
		log.Fatalln("Failed to generate api blueprint:", err)
	}

	err = provide.Provide(*fConfig)
	if err != nil {
		log.Fatalln("Failed to load some services:", err)
	}

	s := server.NewServer()
	err = s.Serve()
	if err != nil {
		log.Fatalln("Failed to serve:", err)
	}
}
