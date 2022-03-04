package middleware

import (
	"fmt"
	"github.com/Aoi-hosizora/ahlib-web/xgin"
	"github.com/Aoi-hosizora/ahlib-web/xgin/headers"
	"github.com/Aoi-hosizora/ahlib/xcolor"
	"github.com/Aoi-hosizora/ahlib/xmodule"
	"github.com/Aoi-hosizora/ahlib/xnumber"
	"github.com/Aoi-hosizora/common_api/internal/pkg/config"
	"github.com/Aoi-hosizora/common_api/internal/pkg/exception"
	"github.com/Aoi-hosizora/common_api/internal/pkg/module/sn"
	"github.com/Aoi-hosizora/common_api/internal/pkg/result"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/juju/ratelimit"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func CorsMiddleware() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:  []string{"*"},
		AllowMethods:  []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:  []string{"Origin", "Content-Type", "Content-Length", "Authorization"},
		ExposeHeaders: []string{"Content-Length"},
	})
}

func getRequestID(c *gin.Context) string {
	return c.Writer.Header().Get(headers.XRequestID)
}

func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		rid := getRequestID(c)
		if rid == "" {
			rid = uuid.New().String()
			c.Header(headers.XRequestID, rid)
		}
	}
}

func LoggerMiddleware() gin.HandlerFunc {
	logger := xmodule.MustGetByName(sn.SLogger).(*logrus.Logger)
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		end := time.Now()

		rid := getRequestID(c)
		options := []xgin.LoggerOption{xgin.WithExtraText(fmt.Sprintf(" | %s", rid)), xgin.WithExtraFieldsV("request_id", rid)}
		xgin.LogResponseToLogrus(logger, c, start, end, options...)
	}
}

func RecoveryMiddleware() gin.HandlerFunc {
	logger := xmodule.MustGetByName(sn.SLogger).(*logrus.Logger)
	const skip = 2
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				errDto, stack := exception.BuildFullErrorDto(err, c, skip) // include request info and trace stack
				xcolor.BrightRed.Printf("\n%s\n\n", stack.String())

				rid := getRequestID(c)
				options := []xgin.LoggerOption{xgin.WithExtraText(fmt.Sprintf(" | %s", rid)), xgin.WithExtraFieldsV("request_id", rid)}
				xgin.LogRecoveryToLogger(logger, err, stack, options...)
				r := result.Error(exception.ServerUnknownError)
				r.Error = errDto
				r.JSON(c)
			}
		}()
		c.Next()
	}
}

func LimiterMiddleware() gin.HandlerFunc {
	cfg := xmodule.MustGetByName(sn.SConfig).(*config.Config).Meta
	limiter := ratelimit.NewBucketWithQuantum(time.Second, cfg.BucketCap, cfg.BucketQua)
	return func(c *gin.Context) {
		available := xnumber.I64toa(limiter.Available())
		capacity := xnumber.I64toa(limiter.Capacity())
		c.Header(headers.XRateLimitRemaining, available)
		c.Header(headers.XRateLimitLimit, capacity)

		if limiter.TakeAvailable(1) == 0 {
			r := gin.H{"remaining": available, "limit": capacity}
			result.Status(http.StatusTooManyRequests).SetData(r).JSON(c)
			c.Abort()
		}
	}
}
