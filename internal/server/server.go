package server

import (
	"context"
	"fmt"
	"github.com/Aoi-hosizora/ahlib-mx/xgin"
	"github.com/Aoi-hosizora/ahlib/xcolor"
	"github.com/Aoi-hosizora/ahlib/xgeneric/xsugar"
	"github.com/Aoi-hosizora/ahlib/xmodule"
	"github.com/Aoi-hosizora/ahlib/xruntime"
	"github.com/Aoi-hosizora/common_api/internal/pkg/apidoc"
	"github.com/Aoi-hosizora/common_api/internal/pkg/config"
	"github.com/Aoi-hosizora/common_api/internal/pkg/module/sn"
	"github.com/Aoi-hosizora/common_api/internal/pkg/result"
	"github.com/Aoi-hosizora/common_api/internal/server/middleware"
	"github.com/Aoi-hosizora/goapidoc"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func init() {
	goapidoc.SetDocument(
		"<placeholder>", "/",
		goapidoc.NewInfo("common_api", "AoiHosizora's common api collection.", "v1.0.0").
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
	engine := xgin.NewEngineSilently(
		xgin.WithDebugPrintRouteFunc(result.PrintRouteFunc(xgin.DefaultColorizedPrintRouteFunc)),
		xgin.WithDefaultWriter(xsugar.If(config.IsDebugMode(), nil, io.Discard)),
		xgin.WithDefaultErrorWriter(xsugar.If(config.IsDebugMode(), nil, io.Discard)),
		xgin.WithRedirectTrailingSlash(true),
		xgin.WithRedirectFixedPath(false),
		xgin.WithRemoveExtraSlash(true),
		xgin.WithHandleMethodNotAllowed(true),
	)

	// middlewares
	engine.Use(middleware.RequestIDMiddleware())
	engine.Use(middleware.LoggerMiddleware())
	engine.Use(middleware.RecoveryMiddleware())
	engine.Use(middleware.LimiterMiddleware())
	engine.Use(middleware.CorsMiddleware())

	// routes
	cfg := xmodule.MustGetByName(sn.SConfig).(*config.Config).Meta
	if cfg.Pprof {
		xgin.WrapPprofSilently(engine)
	}
	if cfg.Swagger {
		xgin.WrapSwagger(engine.Group("/v1/swagger"), apidoc.ReadSwaggerDoc, apidoc.SwaggerOptions()...)
	}
	setupRoutes(engine)

	s := &Server{engine: engine}
	return s, nil
}

func (s *Server) Serve() {
	cfg := xmodule.MustGetByName(sn.SConfig).(*config.Config).Meta
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	server := &http.Server{Addr: addr, Handler: s.engine}

	proxyEnv := xruntime.GetProxyEnv()
	proxyEnv.PrintLog(nil, "[Gin] ")

	terminated := make(chan struct{})
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

	log.Println(xcolor.Bold.Sprintf("[Gin] Listening and serving HTTP on %s", addr))
	fmt.Println()
	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatalln("Failed to serve:", err)
	}
	<-terminated
	log.Println("[Gin] HTTP server is shut down successfully")
}
