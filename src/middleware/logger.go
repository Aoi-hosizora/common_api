package middleware

import (
	"fmt"
	"github.com/Aoi-hosizora/common_api/src/common/logger"
	"github.com/gofiber/fiber"
	"github.com/sirupsen/logrus"
	"time"
)

func LoggerMiddleware() func(c *fiber.Ctx) {
	return func(c *fiber.Ctx) {
		start := time.Now()
		c.Next()
		stop := time.Now()

		method := c.Method()
		path := c.Path()
		latency := stop.Sub(start).String()
		ip := c.IP()

		code := c.Fasthttp.Response.StatusCode()
		length := len(c.Fasthttp.Response.Body())

		entry := logger.Logger.WithFields(logrus.Fields{
			"module":    "fiber",
			"method":    method,
			"path":      path,
			"latency":   latency,
			"code":      code,
			"length":    length,
			"client_ip": ip,
		})
		msg := fmt.Sprintf("[Fiber] %3d | %12s | %15s | %6dB | %-7s %s", code, latency, ip, length, method, path)
		if code >= 500 {
			entry.Error(msg)
		} else if code >= 400 {
			entry.Warn(msg)
		} else {
			entry.Info(msg)
		}
	}
}
