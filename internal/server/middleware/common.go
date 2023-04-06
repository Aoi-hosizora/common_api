package middleware

import (
	"fmt"
	"github.com/Aoi-hosizora/ahlib-mx/xgin"
	"github.com/Aoi-hosizora/ahlib/xcolor"
	"github.com/Aoi-hosizora/ahlib/xconstant/headers"
	"github.com/Aoi-hosizora/ahlib/xmodule"
	"github.com/Aoi-hosizora/ahlib/xnumber"
	"github.com/Aoi-hosizora/ahlib/xreflect"
	"github.com/Aoi-hosizora/common_api/internal/pkg/config"
	"github.com/Aoi-hosizora/common_api/internal/pkg/errno"
	"github.com/Aoi-hosizora/common_api/internal/pkg/module/sn"
	"github.com/Aoi-hosizora/common_api/internal/pkg/result"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/juju/ratelimit"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
	"time"
)

func CorsMiddleware() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Content-Length", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	})
}

func getRequestID(c *gin.Context) string {
	return c.Writer.Header().Get(headers.XRequestID)
}

func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		rid := strings.TrimSpace(getRequestID(c))
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
				errDto, stack := errno.BuildFullErrorDto(err, c, skip) // include request info and trace info
				xcolor.BrightRed.Printf("\n%s\n\n", stack.String())

				rid := getRequestID(c)
				options := []xgin.LoggerOption{xgin.WithExtraText(fmt.Sprintf(" | %s", rid)), xgin.WithExtraFieldsV("request_id", rid)}
				xgin.LogRecoveryToLogrus(logger, err, stack, options...)

				r := result.Error(errno.ServerUnknownError)
				r.Error = errDto
				r.JSON(c)
			}
		}()
		c.Next()
	}
}

func LimiterMiddleware() gin.HandlerFunc {
	cfg := xmodule.MustGetByName(sn.SConfig).(*config.Config).Server
	fillInterval := time.Second * time.Duration(cfg.BucketPeriod)
	limiter := ratelimit.NewBucketWithQuantum(fillInterval, int64(cfg.BucketCap), int64(cfg.BucketQua))
	startTime := xreflect.GetUnexportedField(xreflect.FieldValueOf(limiter, "startTime")).Interface().(time.Time)
	return func(c *gin.Context) {
		available := xnumber.I64toa(limiter.Available())
		capacity := xnumber.I64toa(limiter.Capacity())
		reset := (fillInterval - (time.Now().Sub(startTime) % fillInterval)).String()
		c.Header(headers.XRateLimitRemaining, available)
		c.Header(headers.XRateLimitLimit, capacity)
		c.Header(headers.XRateLimitReset, reset)
		c.Header("X-RateLimit-Policy", fmt.Sprintf("%d;q=%d;w=%d", cfg.BucketCap, cfg.BucketQua, cfg.BucketPeriod))

		if limiter.TakeAvailable(1) == 0 {
			r := gin.H{"remaining": available, "limit": capacity /* always 0 here */, "reset": reset}
			result.Status(http.StatusTooManyRequests).SetData(r).JSON(c)
			c.Abort()
		}
	}
}
