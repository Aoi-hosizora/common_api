package server

import (
	"github.com/Aoi-hosizora/common_api/src/common/result"
	"github.com/gin-gonic/gin"
)

func initRoute(engine *gin.Engine) {
	engine.HandleMethodNotAllowed = true
	engine.NoRoute(func(c *gin.Context) {
		result.Status(404).SetMessage("route not found").JSON(c)
	})
	engine.NoMethod(func(c *gin.Context) {
		result.Status(405).JSON(c)
	})
	engine.GET("/ping", func(c *gin.Context) {
		result.Ok().SetData(&gin.H{"ping": "pong"}).JSON(c)
	})

	// gh := engine.Group("/gh")
	// gh.GET("/users/:name/issues/timeline", controller.GetIssueTimeline)
}
