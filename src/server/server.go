package server

import (
	"fmt"
	"github.com/Aoi-hosizora/ahlib/xdi"
	"github.com/Aoi-hosizora/common_api/docs"
	"github.com/Aoi-hosizora/common_api/src/config"
	"github.com/Aoi-hosizora/common_api/src/middleware"
	"github.com/Aoi-hosizora/common_api/src/provide/sn"
	"github.com/Aoi-hosizora/goapidoc"
	"github.com/DeanThompson/ginpprof"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

func init() {
	goapidoc.SetDocument(
		"localhost:10015", "/",
		goapidoc.NewInfo("gin-n-scaffold", "A go/gin (.NET core style) scaffold template", "1.0").
			WithContact(goapidoc.NewContact("Aoi-hosizora", "https://github.com/Aoi-hosizora", "aoihosizora@hotmail.com")),
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
	engine.Use(middleware.LoggerMiddleware())
	engine.Use(middleware.RecoveryMiddleware())
	engine.Use(middleware.CorsMiddleware())

	// route
	if gin.Mode() == gin.DebugMode {
		ginpprof.Wrap(engine)
	}
	docs.RegisterSwag()
	swaggerUrl := ginSwagger.URL(fmt.Sprintf("http://localhost:%d/swagger/doc.json", cfg.Meta.Port))
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, swaggerUrl))
	initRoute(engine)

	return &Server{engine: engine, config: cfg}
}

func (s *Server) Serve() error {
	addr := fmt.Sprintf("0.0.0.0:%d", s.config.Meta.Port)
	return s.engine.Run(addr)
}
