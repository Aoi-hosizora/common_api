package controller

import (
	"github.com/Aoi-hosizora/common_api/src/common/result"
	"github.com/Aoi-hosizora/goapidoc"
	"github.com/gin-gonic/gin"
)

func init() {
	goapidoc.AddPaths(
		goapidoc.NewPath("GET", "/default", "get common api default message").
			WithTags("Default").
			WithResponses(
				goapidoc.NewResponse(200).WithType("Result"),
			),
	)
}

type DefaultController struct{}

func NewDefaultController() *DefaultController {
	return &DefaultController{}
}

func (d *DefaultController) DefaultMessage(c *gin.Context) {
	result.Ok().SetData(&gin.H{"text": "Here is default group."}).JSON(c)
}
