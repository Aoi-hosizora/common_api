package server

import (
	"fmt"
	"github.com/Aoi-hosizora/common_api/internal/controller"
	"github.com/Aoi-hosizora/common_api/internal/pkg/result"
	"github.com/gin-gonic/gin"
	"strings"
)

func setupRouter(engine *gin.Engine) {
	// meta
	engine.HandleMethodNotAllowed = true
	engine.NoRoute(func(c *gin.Context) {
		result.Status(404).SetMessage(fmt.Sprintf("route '%s' is not found", c.Request.URL.Path)).JSON(c)
	})
	engine.NoMethod(func(c *gin.Context) {
		result.Status(405).SetMessage(fmt.Sprintf("method '%s' is not allowed", strings.ToUpper(c.Request.Method))).JSON(c)
	})
	engine.GET("/ping", func(c *gin.Context) {
		c.JSON(200, &gin.H{"ping": "pong"})
	})
	engine.GET("", func(c *gin.Context) {
		c.JSON(200, &gin.H{"message": "Welcome to AoiHosizora's common-api."})
	})

	// controller

	githubGroup := engine.Group("github")
	{
		githubCtrl := controller.NewGithubController()
		githubGroup.GET("rate_limit", githubCtrl.GetRateLimit)
		githubGroup.GET("users/:name/issues/timeline", githubCtrl.GetIssueTimeline)
		githubGroup.GET("raw", githubCtrl.GetRawPage)
	}

	scutGroup := engine.Group("scut")
	{
		scutCtrl := controller.NewScutController()
		scutGroup.GET("jw", scutCtrl.GetJwItems)
		scutGroup.GET("se", scutCtrl.GetSeItems)
	}
}
