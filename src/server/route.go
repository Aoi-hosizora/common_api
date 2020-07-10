package server

import (
	"github.com/Aoi-hosizora/common_api/src/common/result"
	"github.com/Aoi-hosizora/common_api/src/controller"
	"github.com/gofiber/fiber"
)

func InitRoute(app *fiber.App) {
	app.Get("/ping", func(c *fiber.Ctx) {
		result.Ok().SetData(&result.H{"ping": "pong"}).JSON(c)
	})

	gh := app.Group("/gh")
	gh.Get("/users/:name/issues/event", controller.GetIssueEvents)

	app.All("/*", func(c *fiber.Ctx) {
		result.Status(404).SetMessage("route not found").JSON(c)
	})
}
