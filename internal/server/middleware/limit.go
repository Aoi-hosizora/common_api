package middleware

import (
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

func LimitMiddleware() gin.HandlerFunc {
	cfg := xmodule.MustGetByName(sn.SConfig).(*config.Config).Meta
	limiter := ratelimit.NewBucketWithQuantum(time.Second, cfg.BucketCap, cfg.BucketQua)
	return func(c *gin.Context) {
		c.Writer.Header().Set("X-RateLimit-Remaining", xnumber.I64toa(limiter.Available()))
		c.Writer.Header().Set("X-RateLimit-Limit", xnumber.I64toa(limiter.Capacity()))

		if limiter.TakeAvailable(1) == 0 {
			r := &gin.H{
				"remaining": limiter.Available(),
				"limit":     limiter.Capacity(),
			}
			result.Status(int32(http.StatusTooManyRequests)).SetData(r).JSON(c)
			c.Abort()
		}
	}
}
