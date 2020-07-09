package server

import (
	"fmt"
	"github.com/Aoi-hosizora/common_api/src/config"
	"github.com/Aoi-hosizora/common_api/src/middleware"
	"github.com/gofiber/fiber"
)

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
	InitRoute(app)

	return &Server{App: app}
}

func (s *Server) Serve() error {
	addr := fmt.Sprintf("0.0.0.0:%d", config.Configs.Meta.Port)
	s.Address = addr
	return s.App.Listen(addr)
}
