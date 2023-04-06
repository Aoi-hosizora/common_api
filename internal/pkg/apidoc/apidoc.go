package apidoc

import (
	"fmt"
	"github.com/Aoi-hosizora/ahlib-mx/xgin"
	"github.com/Aoi-hosizora/ahlib/xmodule"
	"github.com/Aoi-hosizora/common_api/internal/pkg/config"
	"github.com/Aoi-hosizora/common_api/internal/pkg/module/sn"
	"github.com/Aoi-hosizora/goapidoc"
	"os"
)

const (
	swaggerDocFilename = "./api/doc.json"
	apibDocFilename    = "./api/doc.apib"
)

func ReadSwaggerDoc() []byte {
	bs, _ := os.ReadFile(swaggerDocFilename)
	return bs
}

func SwaggerOptions() []xgin.SwaggerOption {
	return []xgin.SwaggerOption{
		xgin.WithSwaggerDefaultModelExpandDepth(999),
		xgin.WithSwaggerDisplayRequestDuration(true),
		xgin.WithSwaggerShowExtensions(true),
		xgin.WithSwaggerShowCommonExtensions(true),
	}
}

func UpdateAndSave() error {
	// update host
	cfg := xmodule.MustGetByName(sn.SConfig).(*config.Config).Meta
	host := cfg.DocHost
	if host == "" {
		host = fmt.Sprintf("localhost:%d", cfg.Port)
	}
	goapidoc.SetHost(host)

	// save
	_, err := goapidoc.SaveSwaggerJson(swaggerDocFilename)
	if err != nil {
		return err
	}
	_, err = goapidoc.SaveApib(apibDocFilename)
	if err != nil {
		return err
	}
	return nil
}
