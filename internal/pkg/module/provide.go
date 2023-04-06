package module

import (
	"fmt"
	"github.com/Aoi-hosizora/ahlib/xmodule"
	"github.com/Aoi-hosizora/common_api/internal/pkg/config"
	"github.com/Aoi-hosizora/common_api/internal/pkg/logger"
	"github.com/Aoi-hosizora/common_api/internal/pkg/module/sn"
	"github.com/Aoi-hosizora/common_api/internal/service"
)

func Provide(configPath string) error {
	// ========
	// 1. basic
	// ========

	// *config.Config
	cfg, err := config.Load(configPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}
	xmodule.ProvideByName(sn.SConfig, cfg)

	// *logrus.Logger
	lgr, err := logger.Setup()
	if err != nil {
		return fmt.Errorf("failed to setup logger: %w", err)
	}
	xmodule.ProvideByName(sn.SLogger, lgr)

	// ===========
	// 2. services
	// ===========

	xmodule.ProvideByName(sn.SHttpService, service.NewHttpService())     // *service.HttpService
	xmodule.ProvideByName(sn.SGithubService, service.NewGithubService()) // *service.GithubService
	xmodule.ProvideByName(sn.SScutService, service.NewScutService())     // *service.ScutService

	return nil
}
