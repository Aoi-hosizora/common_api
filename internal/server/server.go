package server

import (
	"context"
	"fmt"
	"github.com/Aoi-hosizora/ahlib-web/xgin"
	"github.com/Aoi-hosizora/ahlib/xmodule"
	"github.com/Aoi-hosizora/ahlib/xruntime"
	"github.com/Aoi-hosizora/common_api/api"
	"github.com/Aoi-hosizora/common_api/internal/pkg/config"
	"github.com/Aoi-hosizora/common_api/internal/pkg/module/sn"
	"github.com/Aoi-hosizora/common_api/internal/server/middleware"
	"github.com/Aoi-hosizora/goapidoc"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func init() {
	goapidoc.SetDocument(
		"localhost:10014", "/",
		goapidoc.NewInfo("common_api", "AoiHosizora's common api collection.", "1.0.0").
			Contact(goapidoc.NewContact("Aoi-hosizora", "https://github.com/Aoi-hosizora", "aoihosizora@hotmail.com")),
	)

	goapidoc.SetOption(goapidoc.NewOption().
		Tags(
			goapidoc.NewTag("Github", "github-controller"),
			goapidoc.NewTag("Scut", "scut-controller"),
		),
	)
}

type Server struct {
	engine *gin.Engine
}

func NewServer() (*Server, error) {
	// server
	restore := xgin.HideDebugLogging()
	engine := gin.New()
	restore()

	// middlewares
	engine.Use(middleware.RequestIDMiddleware())
	engine.Use(middleware.LoggerMiddleware())
	engine.Use(middleware.RecoveryMiddleware())
	engine.Use(middleware.LimiterMiddleware())
	engine.Use(middleware.CorsMiddleware())

	// routes
	xgin.SetPrintRouteFunc(xgin.DefaultColorizedPrintRouteFunc)
	if config.IsDebugMode() {
		restore = xgin.HideDebugPrintRoute()
		xgin.WrapPprof(engine)
		restore()
	}
	cfg := xmodule.MustGetByName(sn.SConfig).(*config.Config)
	if cfg.Meta.Swagger {
		api.RegisterSwagger()
		engine.GET("/swagger/*any", api.SwaggerHandler("doc.json"))
		engine.GET("/swagger", xgin.RedirectHandler(301, "/v1/swagger/index.html"))
	}
	setupRoutes(engine)

	s := &Server{engine: engine}
	return s, nil
}

func (s *Server) Serve() {
	cfg := xmodule.MustGetByName(sn.SConfig).(*config.Config)
	addr := fmt.Sprintf("0.0.0.0:%d", cfg.Meta.Port)
	server := &http.Server{Addr: addr, Handler: s.engine}

	terminated := make(chan interface{})
	go func() {
		defer close(terminated)
		ch := make(chan os.Signal)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
		sig := <-ch
		signal.Stop(ch)
		log.Printf("[Gin] Shutting down due to %s received...", xruntime.SignalName(sig.(syscall.Signal)))

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
		defer cancel()
		err := server.Shutdown(ctx)
		if err != nil {
			log.Fatalln("Failed to shut down:", err)
		}
	}()

	hp, hsp, _ := xruntime.GetProxyEnv()
	if hp != "" {
		log.Printf("[Gin] Using http proxy: %s", hp)
	}
	if hsp != "" {
		log.Printf("[Gin] Using https proxy: %s", hsp)
	}
	log.Printf("[Gin] Listening and serving HTTP on %s", addr)
	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatalln("Failed to serve:", err)
	}
	<-terminated
	log.Println("[Gin] HTTP server is shut down successfully")
}
