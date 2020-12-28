package sn

import (
	"github.com/Aoi-hosizora/ahlib/xdi"
)

const (
	// common
	SConfig xdi.ServiceName = "config" // *config.Config
	SLogger xdi.ServiceName = "logger" // *logrus.Logger

	SHttpService   xdi.ServiceName = "http-service"   // *service.HttpService
	SGithubService xdi.ServiceName = "github-service" // *service.GithubService
	SScutService   xdi.ServiceName = "scut-service"   // *service.ScutService
)
