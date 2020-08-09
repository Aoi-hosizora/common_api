package sn

import (
	"github.com/Aoi-hosizora/ahlib/xdi"
)

const (
	// /src/config/config.go
	SConfig xdi.ServiceName = "config"

	// /src/common/logger/logger.go
	SLogger xdi.ServiceName = "logger"

	// /src/service/*.go
	SHttpService   xdi.ServiceName = "http-service"
	SGithubService xdi.ServiceName = "github-service"
)
