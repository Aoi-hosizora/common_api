package biz

import (
	"time"
)

type GithubIssueSearchResult struct {
	TotalCount        int32              `json:"total_count"`
	IncompleteResults bool               `json:"incomplete_results"`
	Items             []*GithubIssueItem `json:"items"`
}

type GithubIssueItem struct {
	Title    string `json:"title"`
	Number   uint64 `json:"number"`
	HtmlUrl  string `json:"html_url"`
	State    string `json:"state"`
	Comments int32  `json:"comments"`
	Labels   []*struct {
		Name string `json:"name"`
	} `json:"labels"`
	CreatedAt time.Time `json:"created_at"`
}
