package sn

import (
	"github.com/Aoi-hosizora/ahlib/xmodule"
)

const (
	SConfig xmodule.ModuleName = "config" // *config.Config
	SLogger xmodule.ModuleName = "logger" // *logrus.Logger

	SHttpService   xmodule.ModuleName = "http-service"   // *service.HttpService
	SGithubService xmodule.ModuleName = "github-service" // *service.GithubService
	SScutService   xmodule.ModuleName = "scut-service"   // *service.ScutService
)
