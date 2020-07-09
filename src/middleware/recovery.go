package middleware

import (
	"fmt"
	"github.com/Aoi-hosizora/common_api/src/common/exception"
	"github.com/Aoi-hosizora/common_api/src/common/result"
	"github.com/gofiber/fiber"
	"github.com/gofiber/recover"
	"log"
)

func RecoveryMiddleware() func(c *fiber.Ctx) {
	skip := 4
	return recover.New(recover.Config{
		Handler: func(c *fiber.Ctx, err error) {
			fmt.Println()
			log.Println("[Recovery] panic recovered:", err)
			r := result.Error(exception.ServerRecoveryError)
			r.Error = exception.BuildErrorDto(err, skip, c, true)
			r.JSON(c)
		},
		Log: false,
	})
}
