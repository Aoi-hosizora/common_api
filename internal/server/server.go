package server

import (
	"context"
	"fmt"
	"github.com/Aoi-hosizora/ahlib-web/xgin"
	"github.com/Aoi-hosizora/ahlib/xmodule"
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
)

func init() {
	host := os.Getenv("SWAGGER_HOST")
	if host == "" {
		host = "localhost"
	}

	goapidoc.SetDocument(
		host+":10014", "/",
		goapidoc.NewInfo("common_api", "Some small common apis, used for AoiHosizora.\n"+
			"For github, please visit https://github.com/Aoi-hosizora/common_private_api", "1.0").
			Contact(goapidoc.NewContact("Aoi-hosizora", "https://github.com/Aoi-hosizora", "aoihosizora@hotmail.com")),
	)

	goapidoc.SetTags(
		goapidoc.NewTag("Github", "github-controller"),
		goapidoc.NewTag("Scut", "scut-controller"),
	)
}

type Server struct {
	engine *gin.Engine
	config *config.Config
}

func NewServer() (*Server, error) {
	cfg := xmodule.MustGetByName(sn.SConfig).(*config.Config)
	gin.SetMode(cfg.Meta.RunMode)
	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		fmt.Printf("[GIN-debug] %-6s %-30s --> %s (%d handlers)\n", httpMethod, absolutePath, handlerName, nuHandlers)
	}
	xgin.PrintAppRouterRegisterFunc = func(index, count int, method, relativePath, handlerFuncname string, handlersCount int, layerFakePath string) {
		pre := "├─"
		if index == count-1 {
			pre = "└─"
		}
		fmt.Printf("[XGIN]   %2s %-6s ~/%-28s --> %s (%d handlers) ==> ~/%s\n", pre, method, relativePath, handlerFuncname, handlersCount, layerFakePath)
	}

	// server
	engine := gin.New()
	s := &Server{engine: engine, config: cfg}

	// mw
	engine.Use(middleware.RequestIdMiddleware())
	engine.Use(middleware.LoggerMiddleware())
	engine.Use(middleware.RecoveryMiddleware())
	engine.Use(middleware.CorsMiddleware())
	engine.Use(middleware.LimitMiddleware())

	// route
	if gin.Mode() == gin.DebugMode {
		xgin.PprofWrap(engine)
	}
	api.RegisterSwag()
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("doc.json")))
	engine.GET("/swagger", func(c *gin.Context) { c.Redirect(http.StatusPermanentRedirect, "/swagger/index.html") })
	setupRouter(engine)

	return s, nil
}

func (s *Server) Serve() {
	addr := fmt.Sprintf("0.0.0.0:%d", s.config.Meta.Port)
	server := &http.Server{
		Addr:    addr,
		Handler: s.engine,
	}

	closeCh := make(chan int)
	go func() {
		signalCh := make(chan os.Signal)
		signal.Notify(signalCh, syscall.SIGINT)
		<-signalCh
		log.Printf("Shutdown server by SIGINT (0x02)")

		err := server.Shutdown(context.Background())
		if err != nil {
			log.Fatalln("Failed to shutdown HTTP server:", err)
		}
		closeCh <- 0
	}()

	log.Printf("Listening and serving HTTP on %s", addr)
	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatalln("Failed to serve:", err)
	}

	<-closeCh
	log.Println("HTTP server exiting...")
}
