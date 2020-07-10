package server

import (
	"fmt"
	"github.com/Aoi-hosizora/common_api/src/config"
	"github.com/Aoi-hosizora/common_api/src/middleware"
	"github.com/Aoi-hosizora/goapidoc"
	fiberSwagger "github.com/arsmn/fiber-swagger"
	"github.com/gofiber/fiber"
)

func init() {
	goapidoc.SetDocument(
		"localhost:10014", "/",
		goapidoc.NewInfo("common_api", "My common api collection", "1.0").
			WithContact(goapidoc.NewContact("Aoi-hosizora", "https://github.com/Aoi-hosizora/", "aoihosizora@hotmail.com")),
	)
}

type Server struct {
	App     *fiber.App
	Address string
}

func NewServer() *Server {
	app := fiber.New()

	// mw
	app.Use(middleware.RecoveryMiddleware())
	app.Use(middleware.LoggerMiddleware())
	app.Use(middleware.CorsMiddleware())
	app.Use(middleware.CompressionMiddleware())

	// rote
	app.Static("/swagger/doc.json", "./docs/doc.json")
	app.Use("/swagger", fiberSwagger.New(fiberSwagger.Config{
		URL:         fmt.Sprintf("http://localhost:%d/swagger/doc.json", config.Configs.Meta.Port),
		DeepLinking: false,
	}))
	InitRoute(app)

	return &Server{App: app}
}

func (s *Server) Serve() error {
	addr := fmt.Sprintf("0.0.0.0:%d", config.Configs.Meta.Port)
	s.Address = addr
	return s.App.Listen(addr)
}
