package static

const (
	GITHUB_API_URL            = "https://api.github.com/"
	GITHUB_RATE_LIMIT_URL     = "https://api.github.com/rate_limit"
	GITHUB_ISSUE_SEARCH_URL   = "https://api.github.com/search/issues?sort=%s&order=%s&q=involves:%s&page=%d&per_page=%d"
	GITHUB_ISSUE_TIMELINE_URL = "https://api.github.com/repos/%s/%s/issues/%d/timeline?per_page=%d"

	GITHUB_ACCEPT         = "application/vnd.github.v3+json"
	GITHUB_ACCEPT_PREVIEW = "application/vnd.github.mockingbird-preview+json"

	GITHUB_DEFAULT_ISSUE_LIMIT = 20
)
