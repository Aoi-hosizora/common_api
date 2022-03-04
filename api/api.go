package api

import (
	"github.com/Aoi-hosizora/ahlib/xmodule"
	"github.com/Aoi-hosizora/common_api/internal/pkg/config"
	"github.com/Aoi-hosizora/common_api/internal/pkg/module/sn"
	"github.com/Aoi-hosizora/goapidoc"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"github.com/swaggo/swag"
	"os"
)

const (
	SwaggerDocFilename = "./api/spec.json"
)

type swagger struct{}

func (s *swagger) ReadDoc() string {
	f, err := os.ReadFile(SwaggerDocFilename)
	if err != nil {
		return ""
	}
	return string(f)
}

func RegisterSwagger() {
	swag.Register(swag.Name, &swagger{})
}

func SwaggerHandler(docUrl string) gin.HandlerFunc {
	return ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL(docUrl))
}

func UpdateApiDoc() {
	cfg := xmodule.MustGetByName(sn.SConfig).(*config.Config).Meta
	if cfg.DocHost != "" {
		goapidoc.SetHost(cfg.DocHost)
	}
}
