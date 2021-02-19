package logger

import (
	"github.com/Aoi-hosizora/ahlib-more/xlogrus"
	"github.com/Aoi-hosizora/ahlib/xmodule"
	"github.com/Aoi-hosizora/common_api/internal/pkg/config"
	"github.com/Aoi-hosizora/common_api/internal/pkg/module/sn"
	"github.com/sirupsen/logrus"
	"time"
)

func Setup() (*logrus.Logger, error) {
	cfg := xmodule.MustGetByName(sn.SConfig).(*config.Config).Meta

	logger := logrus.New()
	logLevel := logrus.WarnLevel
	if cfg.RunMode == "debug" {
		logLevel = logrus.DebugLevel
	}

	logger.SetLevel(logLevel)
	logger.SetReportCaller(false)
	logger.SetFormatter(&xlogrus.SimpleFormatter{TimestampFormat: time.RFC3339})
	logger.AddHook(xlogrus.NewRotateLogHook(&xlogrus.RotateLogConfig{
		Filename:         cfg.LogName,
		FilenameTimePart: ".%Y%m%d.log",
		LinkFileName:     cfg.LogName + ".log",
		Level:            logLevel,
		Formatter:        &logrus.JSONFormatter{TimestampFormat: time.RFC3339},
		MaxAge:           15 * 24 * time.Hour,
		RotationTime:     24 * time.Hour,
	}))

	return logger, nil
}
