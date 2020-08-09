package logger

import (
	"github.com/Aoi-hosizora/ahlib/xdi"
	"github.com/Aoi-hosizora/ahlib/xlogger"
	"github.com/Aoi-hosizora/common_api/src/config"
	"github.com/Aoi-hosizora/common_api/src/provide/sn"
	"github.com/sirupsen/logrus"
	"time"
)

func Setup() (*logrus.Logger, error) {
	c := xdi.GetByNameForce(sn.SConfig).(*config.Config)

	logger := logrus.New()
	logLevel := logrus.WarnLevel
	if c.Meta.RunMode == "debug" {
		logLevel = logrus.DebugLevel
	}

	logger.SetLevel(logLevel)
	logger.SetReportCaller(false)
	logger.AddHook(xlogger.NewRotateLogHook(&xlogger.RotateLogConfig{
		MaxAge:       15 * 24 * time.Hour,
		RotationTime: 24 * time.Hour,
		Filepath:     c.Meta.LogPath,
		Filename:     c.Meta.LogName,
		Level:        logLevel,
		Formatter:    &logrus.JSONFormatter{TimestampFormat: time.RFC3339},
	}))
	logger.SetFormatter(&xlogger.CustomFormatter{
		ForceColor: true,
	})

	return logger, nil
}
