package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RequestIdMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		rid := c.Writer.Header().Get("X-Request-ID")
		if rid == "" {
			rid = uuid.New().String()
		}

		c.Header("X-Request-ID", rid)
	}
}
