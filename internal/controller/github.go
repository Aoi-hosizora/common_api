package controller

import (
	"encoding/json"
	"github.com/Aoi-hosizora/ahlib/xconstant/headers"
	"github.com/Aoi-hosizora/ahlib/xmodule"
	"github.com/Aoi-hosizora/common_api/internal/model/dto"
	"github.com/Aoi-hosizora/common_api/internal/model/param"
	"github.com/Aoi-hosizora/common_api/internal/pkg/config"
	"github.com/Aoi-hosizora/common_api/internal/pkg/errno"
	"github.com/Aoi-hosizora/common_api/internal/pkg/module/sn"
	"github.com/Aoi-hosizora/common_api/internal/pkg/result"
	"github.com/Aoi-hosizora/common_api/internal/service"
	"github.com/Aoi-hosizora/goapidoc"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"sync"
	"time"
)

func init() {
	goapidoc.AddOperations(
		goapidoc.NewGetOperation("/github/rate_limit", "Get rate limit status for the authenticated user").
			Desc("See https://api.github.com/en/rest/reference/rate-limit").
			Tags("Github").
			AddParams(goapidoc.NewHeaderParam("Authorization", "string", true, "github access token")).
			Responses(goapidoc.NewResponse(200, "string")), // ...

		goapidoc.NewGetOperation("/github/token/{token}/api/{url}", "Request api with given token").
			Tags("Github").
			Deprecated(true).
			AddParams(goapidoc.NewPathParam("token", "string", true, "github access token")).
			AddParams(goapidoc.NewPathParam("url", "string", true, "github api url without api.github.com prefix")).
			Responses(goapidoc.NewResponse(200, "string")), // ...

		goapidoc.NewGetOperation("/github/api/{url}", "Request api with given token").
			Tags("Github").
			AddParams(goapidoc.NewPathParam("url", "string", true, "github api url without api.github.com prefix")).
			AddParams(goapidoc.NewQueryParam("token", "string", true, "github access token")).
			Responses(goapidoc.NewResponse(200, "string")), // ...

		goapidoc.NewGetOperation("/github/repos/{owner}/{repo}/issues", "Get repo simplified issue list").
			Tags("Github").
			AddParams(goapidoc.NewPathParam("owner", "string", true, "owner name")).
			AddParams(goapidoc.NewPathParam("repo", "string", true, "repo name")).
			AddParams(goapidoc.NewQueryParam("page", "integer", false, "current page")).
			AddParams(goapidoc.NewQueryParam("limit", "integer", false, "page size")).
			AddParams(goapidoc.NewHeaderParam("Authorization", "string", true, "github access token")).
			Responses(goapidoc.NewResponse(200, "_Result<GithubIssueItemDto>")),

		goapidoc.NewGetOperation("/github/repos/{owner}/{repo}/issues/search/{q}", "Query repo simplified issue list by title").
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

		goapidoc.NewGetOperation("/github/profile/aoihosizora", "Get Aoi-hosizora user profile with some private fields (default to cache in 1 hour)").
			Tags("Github").
			AddParams(goapidoc.NewQueryParam("force_refresh", "boolean", false, "force to refresh, without using cache")).
			Responses(goapidoc.NewResponse(200, "string")), // ...

		goapidoc.NewGetOperation("/github/contribution/aoihosizora", "Get Aoi-hosizora contribution status (default to cache in 1 hour)").
			Tags("Github").
			AddParams(goapidoc.NewQueryParam("force_refresh", "boolean", false, "force to refresh, without using cache")).
			Responses(goapidoc.NewResponse(200, "string")), // ...
	)
}

type GithubController struct {
	config        *config.Config
	githubService *service.GithubService

	cache *sync.Map
}

func NewGithubController() *GithubController {
	return &GithubController{
		config:        xmodule.MustGetByName(sn.SConfig).(*config.Config),
		githubService: xmodule.MustGetByName(sn.SGithubService).(*service.GithubService),

		cache: &sync.Map{},
	}
}

func (g *GithubController) checkToken(token string) string {
	if strings.TrimSpace(token) == "" {
		return ""
	}
	if strings.HasPrefix(token, "Token ") {
		return token
	}
	return "Token " + token
}

// GetRateLimit GET /github/rate_limit
func (g *GithubController) GetRateLimit(c *gin.Context) *result.Result {
	token := g.checkToken(param.BindToken(c))
	if token == "" {
		return result.Error(errno.RequestParamError)
	}

	rl, err := g.githubService.GetRateLimit(token)
	if err != nil {
		return result.Error(errno.GithubQueryRateLimitError).SetError(err, c)
	}
	c.JSON(http.StatusOK, rl)
	return nil
}

// RequestApiWithToken GET /github/token/:token/api/*url (deprecated)
// RequestApiWithToken GET /github/api/*url?token
func (g *GithubController) RequestApiWithToken(c *gin.Context) *result.Result {
	token := strings.TrimSpace(c.Param("token"))
	if token == "" {
		token = param.BindToken(c)
	}
	token = g.checkToken(token)
	if token == "" {
		return result.Error(errno.RequestParamError)
	}
	url := strings.TrimSpace(c.Param("url"))
	if url == "" {
		return result.Error(errno.RequestParamError)
	}

	bs, statusCode, header, err := g.githubService.RequestApiWithToken(url, token)
	if err != nil {
		return result.Error(errno.GithubQueryApiResponseError).SetError(err, c)
	}
	obj := make(map[string]any)
	if json.Unmarshal(bs, &obj) != nil {
		c.Data(statusCode, header.Get(headers.ContentType), bs)
	} else {
		c.JSON(statusCode, gin.H{"code": statusCode, "data": obj, "headers": header})
	}
	return nil
}

// GetRepoIssues GET /github/repos/:owner/:repo/issues
func (g *GithubController) GetRepoIssues(c *gin.Context) *result.Result {
	owner, repo := c.Param("owner"), c.Param("repo")
	pa := param.BindQueryPage(c, g.config.Github.DefLimit, g.config.Github.MaxLimit)
	token := g.checkToken(param.BindToken(c))
	if token == "" {
		return result.Error(errno.RequestParamError)
	}

	total, items, err := g.githubService.GetRepoIssuesByTitle(owner, repo, pa.Page, pa.Limit, "", token)
	if err != nil {
		return result.Error(errno.GithubQueryRepoIssuesError).SetError(err, c)
	}

	res := dto.BuildGithubIssueItemDtos(items)
	return result.Ok().SetPage(pa.Page, pa.Limit, total, res)
}

// QueryRepoIssuesByTitle GET /github/repos/:owner/:repo/issues/search/:q
func (g *GithubController) QueryRepoIssuesByTitle(c *gin.Context) *result.Result {
	owner, repo := c.Param("owner"), c.Param("repo")
	q := c.Param("q")
	pa := param.BindQueryPage(c, g.config.Github.DefLimit, g.config.Github.MaxLimit)
	token := g.checkToken(param.BindToken(c))
	if token == "" {
		return result.Error(errno.RequestParamError)
	}

	total, items, err := g.githubService.GetRepoIssuesByTitle(owner, repo, pa.Page, pa.Limit, q, token)
	if err != nil {
		return result.Error(errno.GithubQueryRepoIssuesError).SetError(err, c)
	}

	res := dto.BuildGithubIssueItemDtos(items)
	return result.Ok().SetPage(pa.Page, pa.Limit, total, res)
}

// GetIssueTimeline GET /github/users/:owner/issues/timeline?page
func (g *GithubController) GetIssueTimeline(c *gin.Context) *result.Result {
	name := c.Param("owner")
	pa := param.BindQueryPage(c, g.config.Github.DefLimit, g.config.Github.MaxLimit)
	token := g.checkToken(param.BindToken(c))
	if token == "" {
		return result.Error(errno.RequestParamError)
	}

	events, err := g.githubService.GetIssueTimelines(name, pa.Page, token)
	if err != nil {
		return result.Error(errno.GithubQueryIssueTimelineError).SetError(err, c)
	}
	c.JSON(http.StatusOK, events)
	return nil
}

// GetAoiHosizoraUserProfile GET /github/profile/aoihosizora
func (g *GithubController) GetAoiHosizoraUserProfile(c *gin.Context) *result.Result {
	if data, ok := g.cache.Load("profile"); ok && !param.BindBoolQuery(c, "force_refresh", false) {
		return result.Ok().SetCode(200).SetData(data)
	}

	token := g.checkToken(g.config.Github.Token)
	if token == "" {
		return result.Error(errno.RequestParamError)
	}

	data, err := g.githubService.GetAoiHosizoraUserProfile(token)
	if err != nil {
		return result.Error(errno.GithubQueryUserProfileError).SetError(err, c)
	}

	if g.config.Github.CacheDuration > 0 {
		g.cache.Store("profile", data)
		go func() {
			<-time.NewTicker(time.Second * time.Duration(g.config.Github.CacheDuration)).C
			g.cache.Delete("profile")
		}()
	}

	return result.Ok().SetData(data)
}

// GetAoiHosizoraContribution GET /github/contribution/aoihosizora
func (g *GithubController) GetAoiHosizoraContribution(c *gin.Context) *result.Result {
	if data, ok := g.cache.Load("contribution"); ok && !param.BindBoolQuery(c, "force_refresh", false) {
		return result.Ok().SetCode(200).SetData(data)
	}

	token := g.checkToken(g.config.Github.Token)
	if token == "" {
		return result.Error(errno.RequestParamError)
	}

	data, err := g.githubService.GetAoiHosizoraContribution(token)
	if err != nil {
		return result.Error(errno.GithubQueryContributionError).SetError(err, c)
	}
	total := int32(0)
	for _, v := range data {
		total += v
	}
	data["total"] = total

	if g.config.Github.CacheDuration > 0 {
		g.cache.Store("contribution", data)
		go func() {
			<-time.NewTicker(time.Second * time.Duration(g.config.Github.CacheDuration)).C
			g.cache.Delete("contribution")
		}()
	}

	return result.Ok().SetData(data)
}
