package module

import (
	"github.com/Aoi-hosizora/ahlib/xmodule"
	"github.com/Aoi-hosizora/common_api/internal/pkg/config"
	"github.com/Aoi-hosizora/common_api/internal/pkg/logger"
	"github.com/Aoi-hosizora/common_api/internal/pkg/module/sn"
	"github.com/Aoi-hosizora/common_api/internal/service"
	"log"
)

func Provide(configPath string) error {
	// *config.Config
	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatalln("Failed to load config:", err)
	}
	xmodule.ProvideName(sn.SConfig, cfg)

	// *logrus.Logger
	lgr, err := logger.Setup()
	if err != nil {
		log.Fatalln("Failed to setup logger:", err)
	}
	xmodule.ProvideName(sn.SLogger, lgr)

	// ========
	// services
	// ========

	xmodule.ProvideName(sn.SHttpService, service.NewHttpService())     // *service.HttpService
	xmodule.ProvideName(sn.SGithubService, service.NewGithubService()) // *service.GithubService
	xmodule.ProvideName(sn.SScutService, service.NewScutService())     // *service.ScutService

	return nil
}
