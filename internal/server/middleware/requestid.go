package middleware

import (
	"github.com/Aoi-hosizora/ahlib-web/xgin/headers"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		rid := c.Writer.Header().Get(headers.XRequestID)
		if rid == "" {
			rid = uuid.New().String()
			c.Header(headers.XRequestID, rid)
		}
	}
}
