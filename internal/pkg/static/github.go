package static

const (
	GithubRateLimitApi     = "https://api.github.com/rate_limit"
	GithubIssueSearchApi   = "https://api.github.com/search/issues?sort=%s&order=%s&q=involves:%s&page=%d&per_page=%d"
	GithubIssueTimelineApi = "https://api.github.com/repos/%s/%s/issues/%d/timeline?per_page=%d"

	GithubAccept            = "application/vnd.github.v3+json"
	GithubAcceptPreview     = "application/vnd.github.mockingbird-preview+json"
	GithubDefaultIssueLimit = 20
)