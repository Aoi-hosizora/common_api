package middleware

import (
	"github.com/Aoi-hosizora/ahlib-web/xrecovery"
	"github.com/Aoi-hosizora/ahlib/xmodule"
	"github.com/Aoi-hosizora/ahlib/xruntime"
	"github.com/Aoi-hosizora/common_api/internal/pkg/exception"
	"github.com/Aoi-hosizora/common_api/internal/pkg/module/sn"
	"github.com/Aoi-hosizora/common_api/internal/pkg/result"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func RecoveryMiddleware() gin.HandlerFunc {
	logger := xmodule.MustGetByName(sn.SLogger).(*logrus.Logger)
	skip := 2
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				r := result.Error(exception.ServerRecoveryError)
				r.Error = exception.BuildFullErrorDto(err, c)

				rid := c.Writer.Header().Get("X-Request-ID")
				xrecovery.LogToLogrus(logger, err, xruntime.RuntimeTraceStack(skip),
					xrecovery.WithExtraText(rid), xrecovery.WithExtraFieldsV("request_id", rid))
				r.JSON(c)
			}
		}()

		// execute the next handler
		c.Next()
	}
}
