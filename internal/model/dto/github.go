package dto

import (
	"github.com/Aoi-hosizora/common_api/internal/model/object"
	"github.com/Aoi-hosizora/goapidoc"
	"time"
)

func init() {
	goapidoc.AddDefinitions(
		goapidoc.NewDefinition("GithubIssueItemDto", "github issue item").
			Properties(
				goapidoc.NewProperty("title", "string", true, "issue title"),
				goapidoc.NewProperty("number", "integer#int64", true, "issue number"),
				goapidoc.NewProperty("url", "string", true, "issue html url"),
				goapidoc.NewProperty("state", "string", true, "issue state"),
				goapidoc.NewProperty("comments_count", "integer#int32", true, "issue comments count"),
				goapidoc.NewProperty("labels", "string[]", true, "issue labels"),
				goapidoc.NewProperty("create_time", "string", true, "issue create time"),
			),
	)
}

type GithubIssueItemDto struct {
	Title         string   `json:"title"`
	Number        uint64   `json:"number"`
	Url           string   `json:"url"`
	State         string   `json:"state"`
	CommentsCount int32    `json:"comments_count"`
	Labels        []string `json:"labels"`
	CreateTime    string   `json:"create_time"`
}

func BuildGithubIssueItemDto(item *object.GithubIssueItem) *GithubIssueItemDto {
	labels := make([]string, len(item.Labels))
	for idx, label := range item.Labels {
		labels[idx] = label.Name
	}
	return &GithubIssueItemDto{
		Title:         item.Title,
		Number:        item.Number,
		Url:           item.HtmlUrl,
		State:         item.State,
		CommentsCount: item.Comments,
		Labels:        labels,
		CreateTime:    item.CreatedAt.In(time.Local).Format(time.RFC3339),
	}
}

func BuildGithubIssueItemDtos(items []*object.GithubIssueItem) []*GithubIssueItemDto {
	out := make([]*GithubIssueItemDto, len(items))
	for idx, item := range items {
		out[idx] = BuildGithubIssueItemDto(item)
	}
	return out
}
