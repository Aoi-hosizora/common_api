package main

import (
	"flag"
	"github.com/Aoi-hosizora/common_api/src/common/logger"
	"github.com/Aoi-hosizora/common_api/src/config"
	"github.com/Aoi-hosizora/common_api/src/server"
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
	err := config.Load(*fConfig)
	if err != nil {
		log.Fatalln("Failed to load config:", err)
	}
	err = logger.Setup()
	if err != nil {
		log.Fatalln("Failed to setup logger:", err)
	}

	s := server.NewServer()
	err = s.Serve()
	if err != nil {
		log.Fatalf("Failed to listen %s: %v\n", s.Address, err)
	}
}
