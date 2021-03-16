package controller

import (
	"github.com/Aoi-hosizora/ahlib/xmodule"
	"github.com/Aoi-hosizora/common_api/internal/model/dto"
	"github.com/Aoi-hosizora/common_api/internal/pkg/exception"
	"github.com/Aoi-hosizora/common_api/internal/pkg/module/sn"
	"github.com/Aoi-hosizora/common_api/internal/pkg/result"
	"github.com/Aoi-hosizora/common_api/internal/service"
	"github.com/Aoi-hosizora/goapidoc"
	"github.com/gin-gonic/gin"
)

func init() {
	goapidoc.AddRoutePaths(
		goapidoc.NewRoutePath("GET", "/scut/jw/ping", "Ping scut jw").
			Tags("Scut").
			Responses(goapidoc.NewResponse(200, "Result")),

		goapidoc.NewRoutePath("GET", "/scut/jw", "Get scut jw").
			Tags("Scut").
			Responses(goapidoc.NewResponse(200, "_Result<_Page<ScutPostItemDto>>")),

		goapidoc.NewRoutePath("GET", "/scut/se/ping", "Ping scut se").
			Tags("Scut").
			Responses(goapidoc.NewResponse(200, "Result")),

		goapidoc.NewRoutePath("GET", "/scut/se", "Get scut se").
			Tags("Scut").
			Responses(goapidoc.NewResponse(200, "_Result<_Page<ScutPostItemDto>>")),
	)
}

type ScutController struct {
	scutService *service.ScutService
}

func NewScutController() *ScutController {
	return &ScutController{
		scutService: xmodule.MustGetByName(sn.SScutService).(*service.ScutService),
	}
}

// GET /scut/jw/ping
func (s *ScutController) PingJw(c *gin.Context) {
	err := s.scutService.PingJw()
	if err != nil {
		result.Error(exception.PingError).SetError(err, c).JSON(c)
		return
	}
	result.Ok().JSON(c)
}

// GET /scut/jw
func (s *ScutController) GetJwItems(c *gin.Context) {
	items, err := s.scutService.GetJwItems()
	if err != nil {
		result.Error(exception.GetScutJwError).SetError(err, c).JSON(c)
		return
	}

	l := int32(len(items))
	res := dto.BuildScutPostItemDtos(items)
	result.Ok().SetPage(1, l, l, res).JSON(c)
}

// GET /scut/se/ping
func (s *ScutController) PingSe(c *gin.Context) {
	err := s.scutService.PingSe()
	if err != nil {
		result.Error(exception.PingError).SetError(err, c).JSON(c)
		return
	}
	result.Ok().JSON(c)
}

// GET /scut/se
func (s *ScutController) GetSeItems(c *gin.Context) {
	items, err := s.scutService.GetSeItems()
	if err != nil {
		result.Error(exception.GetScutSeError).SetError(err, c).JSON(c)
		return
	}

	l := int32(len(items))
	res := dto.BuildScutPostItemDtos(items)
	result.Ok().SetPage(1, l, l, res).JSON(c)
}
