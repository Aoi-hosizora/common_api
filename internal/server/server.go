package server

import (
	"context"
	"fmt"
	"github.com/Aoi-hosizora/ahlib-web/xgin"
	"github.com/Aoi-hosizora/ahlib/xcolor"
	"github.com/Aoi-hosizora/ahlib/xmodule"
	"github.com/Aoi-hosizora/ahlib/xruntime"
	"github.com/Aoi-hosizora/common_api/api"
	"github.com/Aoi-hosizora/common_api/internal/pkg/config"
	"github.com/Aoi-hosizora/common_api/internal/pkg/module/sn"
	"github.com/Aoi-hosizora/common_api/internal/server/middleware"
	"github.com/Aoi-hosizora/goapidoc"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
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
	if config.IsDebugMode() {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	gin.DebugPrintRouteFunc = func(method, path, handlerName string, numHandlers int) {
		fmt.Printf("[Gin-debug] %s --> %s (%d handlers)\n", xcolor.Blue.Sprintf("%-6s %-30s", method, path), handlerName, numHandlers)
	}

	// server
	engine := xgin.NewEngineWithoutDebugWarning()

	// mw
	engine.Use(middleware.RequestIDMiddleware())
	engine.Use(middleware.LoggerMiddleware())
	engine.Use(middleware.RecoveryMiddleware())
	engine.Use(middleware.LimiterMiddleware())
	engine.Use(middleware.CorsMiddleware())

	// route
	if config.IsDebugMode() {
		restore := xgin.HideDebugPrintRoute()
		xgin.WrapPprof(engine)
		restore()
	}
	cfg := xmodule.MustGetByName(sn.SConfig).(*config.Config)
	if cfg.Meta.Swagger {
		api.RegisterSwagger()
		engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("doc.json")))
		engine.GET("/swagger", func(c *gin.Context) { c.Redirect(http.StatusMovedPermanently, "/swagger/index.html") })
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
		log.Printf("[Gin] Shutting down due to %s received...", xruntime.SignalName(sig.(syscall.Signal)))

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
		defer cancel()
		err := server.Shutdown(ctx)
		if err != nil {
			log.Fatalln("Failed to shut down:", err)
		}
	}()

	hp, hsp, _ := xgin.GetProxyEnv()
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
