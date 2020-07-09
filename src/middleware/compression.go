package middleware

import (
	"github.com/gofiber/compression"
	"github.com/gofiber/fiber"
)

func CompressionMiddleware() func(c *fiber.Ctx) {
	return compression.New()
}
