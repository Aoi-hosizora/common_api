package middleware

import (
	"fmt"
	"github.com/Aoi-hosizora/ahlib-web/xgin"
	"github.com/Aoi-hosizora/ahlib-web/xgin/headers"
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

		rid := c.Writer.Header().Get(headers.XRequestID)
		options := []xgin.LoggerOption{xgin.WithExtraText(fmt.Sprintf(" | %s", rid)), xgin.WithExtraFieldsV("request_id", rid)}
		xgin.LogResponseToLogrus(logger, c, start, end, options...)
	}
}
