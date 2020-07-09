package middleware

import (
	"github.com/Aoi-hosizora/common_api/src/config"
	"github.com/gofiber/cors"
	"github.com/gofiber/fiber"
)

func CorsMiddleware() func(c *fiber.Ctx) {
	origins := []string{"*"}
	if config.Configs.Meta.RunMode == "release" {
		origins = []string{"http://xxx.yyy.zzz", "https://xxx.yyy.zzz"}
	}
	return cors.New(cors.Config{
		AllowOrigins: origins,
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Lines-CheckBox", "Authorization", "Origin"},
	})
}
