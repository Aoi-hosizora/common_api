package middleware

import (
	"github.com/Aoi-hosizora/ahlib-web/xgin"
	"github.com/Aoi-hosizora/ahlib/xdi"
	"github.com/Aoi-hosizora/common_api/src/provide/sn"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"time"
)

func LoggerMiddleware() gin.HandlerFunc {
	lgr := xdi.GetByNameForce(sn.SLogger).(*logrus.Logger)

	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		xgin.LoggerWithLogrus(lgr, start, c)
	}
}
