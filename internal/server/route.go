package server

import (
	"fmt"
	"github.com/Aoi-hosizora/common_api/internal/controller"
	"github.com/Aoi-hosizora/common_api/internal/pkg/result"
	"github.com/gin-gonic/gin"
	"strings"
)

func setupRoutes(engine *gin.Engine) {
	// ===========
	// meta routes
	// ===========

	engine.NoRoute(func(c *gin.Context) {
		msg := fmt.Sprintf("route '%s' is not found", c.Request.URL.Path)
		result.Status(404).SetMessage(msg).JSON(c)
	})
	engine.NoMethod(func(c *gin.Context) {
		msg := fmt.Sprintf("method '%s' is not allowed", strings.ToUpper(c.Request.Method))
		result.Status(405).SetMessage(msg).JSON(c)
	})
	engine.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"ping": "pong"})
	})
	engine.GET("", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Here is AoiHosizora's common api."})
	})

	// ===========
	// controllers
	// ===========

	var (
		githubController = controller.NewGithubController()
		scutController   = controller.NewScutController()
	)

	// ============
	// route groups
	// ============

	// v1 route group
	v1 := result.NewRouteRegisterer(engine) // ATTENTION: no "/v1" prefix

	githubGroup := v1.Group("github")
	githubGroup.GET("rate_limit", githubController.GetRateLimit)
	githubGroup.GET("token/:token/api/*url", githubController.RequestApiWithToken)
	githubGroup.GET("api/*url", githubController.RequestApiWithToken)
	githubGroup.GET("repos/:owner/:repo/issues", githubController.GetRepoIssues)
	githubGroup.GET("repos/:owner/:repo/issues/search/:q", githubController.QueryRepoIssuesByTitle)
	githubGroup.GET("users/:owner/issues/timeline", githubController.GetIssueTimeline)
	githubGroup.GET("profile/aoihosizora", githubController.GetAoiHosizoraUserProfile)

	scutGroup := v1.Group("scut")
	scutGroup.GET("notice/jw", scutController.GetJwNotices)
	scutGroup.GET("notice/se", scutController.GetSeNotices)
	scutGroup.GET("notice/gr", scutController.GetGrNotices)
	scutGroup.GET("notice/gzic", scutController.GetGzicNotices)
}
