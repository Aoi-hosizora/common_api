package middleware

import (
	"github.com/Aoi-hosizora/ahlib-web/xgin"
	"github.com/Aoi-hosizora/ahlib-web/xgin/headers"
	"github.com/Aoi-hosizora/ahlib/xcolor"
	"github.com/Aoi-hosizora/ahlib/xmodule"
	"github.com/Aoi-hosizora/common_api/internal/pkg/exception"
	"github.com/Aoi-hosizora/common_api/internal/pkg/module/sn"
	"github.com/Aoi-hosizora/common_api/internal/pkg/result"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func RecoveryMiddleware() gin.HandlerFunc {
	logger := xmodule.MustGetByName(sn.SLogger).(*logrus.Logger)
	const skip = 3
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				errDto, stack := exception.BuildFullErrorDto(err, c, skip) // include request info and trace stack
				xcolor.BrightRed.Printf("\n%s\n\n", stack.String())

				rid := c.Writer.Header().Get(headers.XRequestID)
				options := []xgin.LoggerOption{xgin.WithExtraText(rid), xgin.WithExtraFieldsV("request_id", rid)}
				xgin.LogRecoveryToLogger(logger, err, stack, options...)

				r := result.Error(exception.ServerUnknownError)
				r.Error = errDto
				r.JSON(c)
			}
		}()
		c.Next()
	}
}
