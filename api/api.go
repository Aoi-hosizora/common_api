package api

import (
	"github.com/Aoi-hosizora/ahlib/xmodule"
	"github.com/Aoi-hosizora/common_api/internal/pkg/config"
	"github.com/Aoi-hosizora/common_api/internal/pkg/module/sn"
	"github.com/Aoi-hosizora/goapidoc"
	"github.com/swaggo/swag"
	"io/ioutil"
)

const (
	SwaggerDocFilename = "./api/spec.json"
)

type swagger struct{}

func (s *swagger) ReadDoc() string {
	f, err := ioutil.ReadFile(SwaggerDocFilename)
	if err != nil {
		return ""
	}
	return string(f)
}

func RegisterSwagger() {
	swag.Register(swag.Name, &swagger{})
}

func UpdateApiDoc() {
	cfg := xmodule.MustGetByName(sn.SConfig).(*config.Config).Meta
	if cfg.Host != "" {
		goapidoc.SetHost(cfg.Host)
	}
}
