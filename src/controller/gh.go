package controller

import (
	"github.com/Aoi-hosizora/common_api/src/common/exception"
	"github.com/Aoi-hosizora/common_api/src/common/result"
	"github.com/Aoi-hosizora/common_api/src/service"
	"github.com/Aoi-hosizora/goapidoc"
	"github.com/gofiber/fiber"
	"strconv"
)

func init() {
	goapidoc.AddPaths(
		goapidoc.NewPath("GET", "/gh/users/{name}/issues/timeline", "get github user issues timeline (event)").
			WithDescription("Fixed field: id?, node_id?, event(enum), actor(User), commit_id?, commit_url?, created_at(time), repo(string), number(integer), involve(string)").
			WithTags("Github").
			WithParams(
				goapidoc.NewPathParam("name", "string", true, "github username"),
				goapidoc.NewQueryParam("page", "integer#int32", false, "query page"),
				goapidoc.NewHeaderParam("Authorization", "string", false, "github token, format: Token xxx"),
			).
			WithResponses(
				goapidoc.NewResponse(200).WithType("string[]"),
			),
	)
}

// /gh/users/:name/issues/timeline?page
func GetIssueTimeline(c *fiber.Ctx) {
	name := c.Params("name")
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil || page <= 0 {
		page = 1
	}
	auth := c.Get("Authorization")

	events, err := service.GetIssueEvents(name, int32(page), auth)
	if err != nil {
		result.Error(exception.GetGithubError).SetError(err, c).JSON(c)
	} else {
		_ = c.JSON(events)
	}
}
