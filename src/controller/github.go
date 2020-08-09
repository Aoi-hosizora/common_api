package controller

import (
	"github.com/Aoi-hosizora/ahlib/xdi"
	"github.com/Aoi-hosizora/common_api/src/common/exception"
	"github.com/Aoi-hosizora/common_api/src/common/result"
	"github.com/Aoi-hosizora/common_api/src/provide/sn"
	"github.com/Aoi-hosizora/common_api/src/service"
	"github.com/Aoi-hosizora/goapidoc"
	"github.com/gin-gonic/gin"
	"strconv"
)

func init() {
	goapidoc.AddPaths(
		goapidoc.NewPath("GET", "/github/users/{name}/issues/timeline", "get github user issues timeline (event)").
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
	)
}

type GithubController struct {
	githubService *service.GithubService
}

func NewGithubService() *GithubController {
	return &GithubController{
		githubService: xdi.GetByNameForce(sn.SGithubService).(*service.GithubService),
	}
}

// /gh/users/:name/issues/timeline?page
func (g *GithubController) GetIssueTimeline(c *gin.Context) {
	name := c.Param("name")
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil || page <= 0 {
		page = 1
	}
	auth := c.GetHeader("Authorization")
	if auth == "" {
		auth = c.DefaultQuery("Authorization", "")
	}
	if auth == "" {
		result.Error(exception.RequestParamError).JSON(c)
		return
	}

	events, err := g.githubService.GetIssueEvents(name, int32(page), auth)
	if err != nil {
		result.Error(exception.GetGithubError).SetError(err, c).JSON(c)
		return
	}

	result.Ok().SetData(events).JSON(c)
}
