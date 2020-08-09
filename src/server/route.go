package server

import (
	"github.com/Aoi-hosizora/common_api/src/common/result"
	"github.com/Aoi-hosizora/common_api/src/controller"
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
	engine.GET("/", func(c *gin.Context) {
		result.Ok().SetData(&gin.H{"text": "Here is AoiHosizora' common api."}).JSON(c)
	})
	engine.GET("/ping", func(c *gin.Context) {
		result.Ok().SetData(&gin.H{"ping": "pong"}).JSON(c)
	})

	// /default
	{
		defaultController := controller.NewDefaultController()
		def := engine.Group("/default")
		def.GET("", defaultController.DefaultMessage)
	}

	// /github
	{
		githubController := controller.NewGithubController()
		github := engine.Group("/github")
		github.GET("/rate_limit", githubController.GetRateLimit)
		github.GET("/users/:name/issues/timeline", githubController.GetIssueTimeline)
	}
}
