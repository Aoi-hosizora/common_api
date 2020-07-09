package main

import (
	"github.com/gofiber/fiber"
	"log"
)

func main() {
	app := fiber.New()

	app.Get("/ping", func(c *fiber.Ctx) {
		_ = c.JSON("pong")
	})

	err := app.Listen(10014)
	if err != nil {
		log.Fatalln(err)
	}
}
