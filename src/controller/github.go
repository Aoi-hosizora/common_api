package controller

import (
	"github.com/Aoi-hosizora/ahlib/xdi"
	"github.com/Aoi-hosizora/common_api/src/common/exception"
	"github.com/Aoi-hosizora/common_api/src/common/result"
	"github.com/Aoi-hosizora/common_api/src/provide/sn"
	"github.com/Aoi-hosizora/common_api/src/service"
	"github.com/Aoi-hosizora/goapidoc"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func init() {
	goapidoc.AddRoutePaths(
		goapidoc.NewRoutePath("GET", "/github/rate_limit", "Get rate limit status for the authenticated user").
			Desc("See https://docs.github.com/en/rest/reference/rate-limit").
			Tags("Github").
			Params(
				goapidoc.NewHeaderParam("Authorization", "string", true, "github token, format: Token xxx"),
			).
			Responses(
				goapidoc.NewResponse(200, "string"),
			),

		goapidoc.NewRoutePath("GET", "/github/users/{name}/issues/timeline", "Get github user issues timeline (event)").
			Desc("Fixed field: id?, node_id?, event(enum), actor(User), commit_id?, commit_url?, created_at(time), repo(string), number(integer), involve(string)").
			Tags("Github").
			Params(
				goapidoc.NewPathParam("name", "string", true, "github username"),
				goapidoc.NewQueryParam("page", "integer#int32", false, "query page"),
				goapidoc.NewHeaderParam("Authorization", "string", true, "github token, format: Token xxx"),
			).
			Responses(
				goapidoc.NewResponse(200, "string[]"),
			),

		goapidoc.NewRoutePath("GET", "/github/raw", "Get raw page without authentication").
			Tags("Github").
			Params(
				goapidoc.NewQueryParam("page", "string", true, "Github url without github.com prefix"),
			).
			Responses(
				goapidoc.NewResponse(200, "string"),
			),
	)
}

type GithubController struct {
	githubService *service.GithubService
}

func NewGithubController() *GithubController {
	return &GithubController{
		githubService: xdi.GetByNameForce(sn.SGithubService).(*service.GithubService),
	}
}

// GET /github/rate_limit?page
func (g *GithubController) GetRateLimit(c *gin.Context) {
	auth := c.GetHeader("Authorization")
	if auth == "" {
		auth = c.DefaultQuery("Authorization", "")
		if auth == "" {
			result.Error(exception.RequestParamError).JSON(c)
			return
		}
	}

	core, err := g.githubService.GetRateLimit(auth)
	if err != nil {
		result.Error(exception.GetGithubError).SetError(err, c).JSON(c)
		return
	}

	c.JSON(http.StatusOK, core)
}

// GET /github/users/:name/issues/timeline?page
func (g *GithubController) GetIssueTimeline(c *gin.Context) {
	name := c.Param("name")
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil || page <= 0 {
		page = 1
	}
	auth := c.GetHeader("Authorization")
	if auth == "" {
		auth = c.DefaultQuery("Authorization", "")
		if auth == "" {
			result.Error(exception.RequestParamError).JSON(c)
			return
		}
	}

	events, err := g.githubService.GetIssueEvents(name, int32(page), auth)
	if err != nil {
		result.Error(exception.GetGithubError).SetError(err, c).JSON(c)
		return
	}

	c.JSON(http.StatusOK, events)
}

// GET /github/raw?page
func (g *GithubController) GetRawPage(c *gin.Context) {
	page := c.DefaultQuery("page", "")
	if page == "" {
		result.Error(exception.RequestParamError).JSON(c)
		return
	}

	html, err := g.githubService.GetRawPage(page)
	if err != nil {
		result.Error(exception.GetGithubError).SetError(err, c).JSON(c)
		return
	}

	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, html)
}
