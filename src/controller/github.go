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
	goapidoc.AddPaths(
		goapidoc.NewPath("GET", "/github/rate_limit", "Get rate limit status for the authenticated user").
			WithDescription("See https://docs.github.com/en/rest/reference/rate-limit").
			WithTags("Github").
			WithParams(
				goapidoc.NewHeaderParam("Authorization", "string", true, "github token, format: Token xxx"),
			).
			WithResponses(
				goapidoc.NewResponse(200).WithType("string"),
			),

		goapidoc.NewPath("GET", "/github/users/{name}/issues/timeline", "Get github user issues timeline (event)").
			WithDescription("Fixed field: id?, node_id?, event(enum), actor(User), commit_id?, commit_url?, created_at(time), repo(string), number(integer), involve(string)").
			WithTags("Github").
			WithParams(
				goapidoc.NewPathParam("name", "string", true, "github username"),
				goapidoc.NewQueryParam("page", "integer#int32", false, "query page"),
				goapidoc.NewHeaderParam("Authorization", "string", true, "github token, format: Token xxx"),
			).
			WithResponses(
				goapidoc.NewResponse(200).WithType("string[]"),
			),

		goapidoc.NewPath("GET", "/github/raw", "Get raw page without authentication").
			WithTags("Github").
			WithParams(
				goapidoc.NewQueryParam("page", "string", true, "Github url without github.com prefix"),
			).
			WithResponses(
				goapidoc.NewResponse(200).WithType("string"),
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

// /github/rate_limit?page
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

// /github/users/:name/issues/timeline?page
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

// /github/raw?page
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
