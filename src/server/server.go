package server

import (
	"fmt"
	"github.com/Aoi-hosizora/ahlib-web/xgin"
	"github.com/Aoi-hosizora/ahlib/xdi"
	"github.com/Aoi-hosizora/common_api/docs"
	"github.com/Aoi-hosizora/common_api/src/config"
	"github.com/Aoi-hosizora/common_api/src/middleware"
	"github.com/Aoi-hosizora/common_api/src/provide/sn"
	"github.com/Aoi-hosizora/goapidoc"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"net/http"
)

func init() {
	goapidoc.SetDocument(
		"localhost:10014", "/",
		goapidoc.NewInfo("common_api", "My common api collection", "1.0").
			Contact(goapidoc.NewContact("Aoi-hosizora", "https://github.com/Aoi-hosizora", "aoihosizora@hotmail.com")),
	)

	goapidoc.SetTags(
		goapidoc.NewTag("Github", "github-controller"),
	)
}

type Server struct {
	engine *gin.Engine
	config *config.Config
}

func NewServer() *Server {
	cfg := xdi.GetByNameForce(sn.SConfig).(*config.Config)
	gin.SetMode(cfg.Meta.RunMode)
	engine := gin.New()

	// mw
	engine.Use(middleware.RequestIdMiddleware())
	engine.Use(middleware.LoggerMiddleware())
	engine.Use(middleware.RecoveryMiddleware())
	engine.Use(middleware.CorsMiddleware())

	// route
	if gin.Mode() == gin.DebugMode {
		xgin.PprofWrap(engine)
	}
	docs.RegisterSwag()
	engine.GET("/v1/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("doc.json")))
	engine.GET("/v1/swagger", func(c *gin.Context) { c.Redirect(http.StatusPermanentRedirect, "/v1/swagger/index.html") })
	initRoute(engine)

	return &Server{engine: engine, config: cfg}
}

func (s *Server) Serve() error {
	addr := fmt.Sprintf("0.0.0.0:%d", s.config.Meta.Port)
	return s.engine.Run(addr)
}
