package logger

import (
	"github.com/Aoi-hosizora/ahlib/xdi"
	"github.com/Aoi-hosizora/common_api/src/provide/sn"
	"github.com/sirupsen/logrus"
)

func LogGhUrl(url string) {
	logger := xdi.GetByNameForce(sn.SLogger).(*logrus.Logger)
	logger.Infof("[Github] %s", url)
}
