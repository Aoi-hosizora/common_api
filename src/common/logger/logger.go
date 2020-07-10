package logger

import (
	"github.com/Aoi-hosizora/ahlib/xlogger"
	"github.com/Aoi-hosizora/common_api/src/config"
	"github.com/sirupsen/logrus"
	"github.com/snowzach/rotatefilehook"
	"runtime"
	"time"
)

var Logger *logrus.Logger

func Setup() error {
	logger := logrus.New()
	logLevel := logrus.WarnLevel
	if config.Configs.Meta.RunMode == "debug" {
		logLevel = logrus.DebugLevel
	}

	fileHook, err := rotatefilehook.NewRotateFileHook(rotatefilehook.RotateFileConfig{
		Filename:   config.Configs.Meta.LogPath,
		MaxSize:    50,
		MaxBackups: 3,
		MaxAge:     30,
		Level:      logLevel,
		Formatter:  &logrus.JSONFormatter{TimestampFormat: time.RFC3339},
	})
	if err != nil {
		return err
	}

	logger.SetLevel(logLevel)
	logger.SetReportCaller(true)
	logger.AddHook(fileHook)
	logger.SetFormatter(&xlogger.CustomerFormatter{
		ForceColor: true,
		RuntimeCaller: func(f *runtime.Frame) (string, string) {
			return "", ""
		},
	})

	Logger = logger
	return nil
}

func LogGhUrl(url string) {
	Logger.Infof("[GITHUB] %s", url)
}
