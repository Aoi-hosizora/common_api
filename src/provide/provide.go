package provide

import (
	"github.com/Aoi-hosizora/ahlib/xdi"
	"github.com/Aoi-hosizora/common_api/src/common/logger"
	"github.com/Aoi-hosizora/common_api/src/config"
	"github.com/Aoi-hosizora/common_api/src/model/profile"
	"github.com/Aoi-hosizora/common_api/src/provide/sn"
	"github.com/Aoi-hosizora/common_api/src/service"
	"log"
)

func Provide(configPath string) error {
	// /src/config/config.go
	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatalln("Failed to load config:", err)
	}
	xdi.ProvideName(sn.SConfig, cfg)

	// /src/common/logger/logger.go
	lgr, err := logger.Setup()
	if err != nil {
		log.Fatalln("Failed to setup logger:", err)
	}
	xdi.ProvideName(sn.SLogger, lgr)

	// /src/service/*.go
	xdi.ProvideName(sn.SHttpService, service.NewHttpService())
	xdi.ProvideName(sn.SGithubService, service.NewGithubService())

	// /src/model/profile/{entity.go, property.go}
	profile.BuildEntityMappers()
	profile.BuildPropertyMappers()

	return nil
}
