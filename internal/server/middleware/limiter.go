package middleware

import (
	"github.com/Aoi-hosizora/ahlib-web/xgin/headers"
	"github.com/Aoi-hosizora/ahlib/xmodule"
	"github.com/Aoi-hosizora/ahlib/xnumber"
	"github.com/Aoi-hosizora/common_api/internal/pkg/config"
	"github.com/Aoi-hosizora/common_api/internal/pkg/module/sn"
	"github.com/Aoi-hosizora/common_api/internal/pkg/result"
	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
	"net/http"
	"time"
)

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
