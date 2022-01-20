package controller

import (
	"github.com/Aoi-hosizora/ahlib/xmodule"
	"github.com/Aoi-hosizora/ahlib/xnumber"
	"github.com/Aoi-hosizora/common_api/internal/pkg/exception"
	"github.com/Aoi-hosizora/common_api/internal/pkg/module/sn"
	"github.com/Aoi-hosizora/common_api/internal/pkg/result"
	"github.com/Aoi-hosizora/common_api/internal/service"
	"github.com/Aoi-hosizora/goapidoc"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func init() {
	goapidoc.AddOperations(
		goapidoc.NewGetOperation("/github/rate_limit", "Get rate limit status for the authenticated user").
			Desc("See https://api.github.com/en/rest/reference/rate-limit").
			Tags("Github").
			AddParams(goapidoc.NewHeaderParam("Authorization", "string", true, "github token, format: Token xxx")).
			Responses(goapidoc.NewResponse(200, "string")), // ...

		goapidoc.NewGetOperation("/github/users/{name}/issues/timeline", "Get github user issues timeline (event)").
			Desc("Fixed field: id?, node_id?, event(enum), actor(User), commit_id?, commit_url?, created_at(time), repo(string), number(integer), involve(string)").
			Tags("Github").
			AddParams(goapidoc.NewPathParam("name", "string", true, "github username")).
			AddParams(goapidoc.NewQueryParam("page", "integer#int32", false, "query page")).
			AddParams(goapidoc.NewHeaderParam("Authorization", "string", true, "github token, format: Token xxx")).
			Responses(goapidoc.NewResponse(200, "string[]")), // ...
	)
}

type GithubController struct {
	githubService *service.GithubService
}

func NewGithubController() *GithubController {
	return &GithubController{
		githubService: xmodule.MustGetByName(sn.SGithubService).(*service.GithubService),
	}
}

func (g *GithubController) token(c *gin.Context) string {
	auth := strings.TrimSpace(c.GetHeader("Authorization"))
	if auth == "" {
		auth = strings.TrimSpace(c.DefaultQuery("token", ""))
	}
	return auth
}

// GetRateLimit GET /github/rate_limit
func (g *GithubController) GetRateLimit(c *gin.Context) {
	auth := g.token(c)
	if auth == "" {
		result.Error(exception.RequestParamError).JSON(c)
		return
	}

	rl, err := g.githubService.GetRateLimit(auth)
	if err != nil {
		result.Error(exception.GithubQueryRateLimitError).SetError(err, c).JSON(c)
		return
	}
	c.JSON(http.StatusOK, rl)
}

// GetIssueTimeline GET /github/users/:name/issues/timeline?page
func (g *GithubController) GetIssueTimeline(c *gin.Context) {
	name := c.Param("name")
	page, err := xnumber.Atoi32(c.Query("page"))
	if err != nil || page <= 0 {
		page = 1
	}
	auth := g.token(c)
	if auth == "" {
		result.Error(exception.RequestParamError).JSON(c)
		return
	}

	events, err := g.githubService.GetIssueTimelines(name, page, auth)
	if err != nil {
		result.Error(exception.GithubQueryIssueTimelineError).SetError(err, c).JSON(c)
		return
	}
	c.JSON(http.StatusOK, events)
}
