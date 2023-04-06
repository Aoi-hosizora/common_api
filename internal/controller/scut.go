package controller

import (
	"github.com/Aoi-hosizora/ahlib/xmodule"
	"github.com/Aoi-hosizora/common_api/internal/model/dto"
	"github.com/Aoi-hosizora/common_api/internal/pkg/errno"
	"github.com/Aoi-hosizora/common_api/internal/pkg/module/sn"
	"github.com/Aoi-hosizora/common_api/internal/pkg/result"
	"github.com/Aoi-hosizora/common_api/internal/service"
	"github.com/Aoi-hosizora/goapidoc"
	"github.com/gin-gonic/gin"
)

func init() {
	goapidoc.AddOperations(
		goapidoc.NewGetOperation("/scut/notice/jw", "Get scut jw notices").
			Tags("Scut").
			Responses(goapidoc.NewResponse(200, "_Result<_Page<ScutNoticeItemDto>>")),

		goapidoc.NewGetOperation("/scut/notice/se", "Get scut se notices").
			Tags("Scut").
			Responses(goapidoc.NewResponse(200, "_Result<_Page<ScutNoticeItemDto>>")),

		goapidoc.NewGetOperation("/scut/notice/gr", "Get scut gr notices").
			Tags("Scut").
			Responses(goapidoc.NewResponse(200, "_Result<_Page<ScutNoticeItemDto>>")),

		goapidoc.NewGetOperation("/scut/notice/gzic", "Get scut gzic notices").
			Tags("Scut").
			Responses(goapidoc.NewResponse(200, "_Result<_Page<ScutNoticeItemDto>>")),
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

// GetJwNotices GET /scut/notice/jw
func (s *ScutController) GetJwNotices(c *gin.Context) *result.Result {
	items, err := s.scutService.GetJwNotices()
	if err != nil {
		return result.Error(errno.ScutQueryJwNoticesError).SetError(err, c)
	}

	l := uint32(len(items))
	res := dto.BuildScutNoticeItemDtos(items)
	return result.Ok().SetPage(1, l, l, res)
}

// GetSeNotices GET /scut/notice/se
func (s *ScutController) GetSeNotices(c *gin.Context) *result.Result {
	items, err := s.scutService.GetSeNotices()
	if err != nil {
		return result.Error(errno.ScutQuerySeNoticesError).SetError(err, c)
	}

	l := uint32(len(items))
	res := dto.BuildScutNoticeItemDtos(items)
	return result.Ok().SetPage(1, l, l, res)
}

// GetGrNotices GET /scut/notice/gr
func (s *ScutController) GetGrNotices(c *gin.Context) *result.Result {
	items, err := s.scutService.GetGrNotices()
	if err != nil {
		return result.Error(errno.ScutQueryGrNoticesError).SetError(err, c)
	}

	l := uint32(len(items))
	res := dto.BuildScutNoticeItemDtos(items)
	return result.Ok().SetPage(1, l, l, res)
}

// GetGzicNotices GET /scut/notice/gzic
func (s *ScutController) GetGzicNotices(c *gin.Context) *result.Result {
	items, err := s.scutService.GetGzicNotices()
	if err != nil {
		return result.Error(errno.ScutQueryGzicNoticesError).SetError(err, c)
	}

	l := uint32(len(items))
	res := dto.BuildScutNoticeItemDtos(items)
	return result.Ok().SetPage(1, l, l, res)
}
