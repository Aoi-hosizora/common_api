package controller

import (
	"github.com/Aoi-hosizora/ahlib/xmodule"
	"github.com/Aoi-hosizora/ahlib/xnumber"
	"github.com/Aoi-hosizora/common_api/internal/model/dto"
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
			AddParams(goapidoc.NewHeaderParam("Authorization", "string", true, "github access token")).
			Responses(goapidoc.NewResponse(200, "string")), // ...

		goapidoc.NewGetOperation("/github/repos/{owner}/{repo}", "Get repo simplified issue list").
			Tags("Github").
			AddParams(goapidoc.NewPathParam("owner", "string", true, "owner name")).
			AddParams(goapidoc.NewPathParam("repo", "string", true, "repo name")).
			AddParams(goapidoc.NewQueryParam("page", "integer", false, "current page")).
			AddParams(goapidoc.NewQueryParam("limit", "integer", false, "page size")).
			AddParams(goapidoc.NewHeaderParam("Authorization", "string", true, "github access token")).
			Responses(goapidoc.NewResponse(200, "_Result<GithubIssueItemDto>")),

		goapidoc.NewGetOperation("/github/repos/{owner}/{repo}/search/{q}", "Query repo simplified issue list by title").
			Tags("Github").
			AddParams(goapidoc.NewPathParam("owner", "string", true, "owner name")).
			AddParams(goapidoc.NewPathParam("repo", "string", true, "repo name")).
			AddParams(goapidoc.NewPathParam("q", "string", true, "issue title")).
			AddParams(goapidoc.NewQueryParam("page", "integer", false, "current page")).
			AddParams(goapidoc.NewQueryParam("limit", "integer", false, "page size")).
			AddParams(goapidoc.NewHeaderParam("Authorization", "string", true, "github access token")).
			Responses(goapidoc.NewResponse(200, "_Result<GithubIssueItemDto>")),

		goapidoc.NewGetOperation("/github/users/{owner}/issues/timeline", "Get user issues timeline (event)").
			Desc("Fixed field: id?, node_id?, event(enum), actor(User), commit_id?, commit_url?, created_at(time), repo(string), number(integer), involve(string)").
			Tags("Github").
			AddParams(goapidoc.NewPathParam("owner", "string", true, "owner name")).
			AddParams(goapidoc.NewQueryParam("page", "integer#int32", false, "query page")).
			AddParams(goapidoc.NewHeaderParam("Authorization", "string", true, "github access token")).
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
	return "Token " + strings.TrimSpace(strings.TrimPrefix(auth, "Token"))
}

func (g *GithubController) bindPage(c *gin.Context, defLimit, maxLimit int32) (page int32, limit int32) {
	page, err := xnumber.Atoi32(c.Query("page"))
	if err != nil || page <= 0 {
		page = 1
	}

	limit, err = xnumber.Atoi32(c.Query("limit"))
	if err != nil || limit <= 0 {
		limit = defLimit
	} else if limit > maxLimit {
		limit = maxLimit
	}

	return page, limit
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

// GetRepoIssues GET /github/repos/:owner/:repo
func (g *GithubController) GetRepoIssues(c *gin.Context) {
	owner := c.Param("owner")
	repo := c.Param("repo")
	page, limit := g.bindPage(c, 20, 50)
	auth := g.token(c)
	if auth == "" {
		result.Error(exception.RequestParamError).JSON(c)
		return
	}

	total, items, err := g.githubService.GetRepoIssuesByTitle(owner, repo, page, limit, "", auth)
	if err != nil {
		result.Error(exception.GithubQueryRepoIssuesError).SetError(err, c).JSON(c)
		return
	}

	res := dto.BuildGithubIssueItemDtos(items)
	result.Ok().SetPage(page, limit, total, res).JSON(c)
}

// QueryRepoIssuesByTitle GET /github/repos/:owner/:repo/search/:q
func (g *GithubController) QueryRepoIssuesByTitle(c *gin.Context) {
	owner := c.Param("owner")
	repo := c.Param("repo")
	q := c.Param("q")
	page, limit := g.bindPage(c, 20, 50)
	auth := g.token(c)
	if auth == "" {
		result.Error(exception.RequestParamError).JSON(c)
		return
	}

	total, items, err := g.githubService.GetRepoIssuesByTitle(owner, repo, page, limit, q, auth)
	if err != nil {
		result.Error(exception.GithubQueryRepoIssuesError).SetError(err, c).JSON(c)
		return
	}

	res := dto.BuildGithubIssueItemDtos(items)
	result.Ok().SetPage(page, limit, total, res).JSON(c)
}

// GetIssueTimeline GET /github/users/:owner/issues/timeline?page
func (g *GithubController) GetIssueTimeline(c *gin.Context) {
	name := c.Param("owner")
	page, _ := g.bindPage(c, 0, 0)
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
