package middleware

import (
	"github.com/Aoi-hosizora/ahlib-web/xgin"
	"github.com/Aoi-hosizora/ahlib/xmodule"
	"github.com/Aoi-hosizora/common_api/internal/pkg/module/sn"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"time"
)

func LoggerMiddleware() gin.HandlerFunc {
	logger := xmodule.MustGetByName(sn.SLogger).(*logrus.Logger)
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		end := time.Now()

		rid := c.Writer.Header().Get("X-Request-ID")
		xgin.LogToLogrus(logger, c, start, end,
			xgin.WithExtraText(rid), xgin.WithExtraFieldsV("request_id", rid))
	}
}
