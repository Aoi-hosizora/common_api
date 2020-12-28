package controller

import (
	"github.com/Aoi-hosizora/ahlib/xdi"
	"github.com/Aoi-hosizora/common_api/src/provide/sn"
	"github.com/Aoi-hosizora/common_api/src/service"
)

type ScutController struct {
	scutService *service.ScutService
}

func NewScutController() *ScutController {
	return &ScutController{
		scutService: xdi.GetByNameForce(sn.SScutService).(*service.ScutService),
	}
}
