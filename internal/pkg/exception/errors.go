package exception

var (
	errno4xx = int32(40000) - 1
	errno5xx = int32(50000) - 1
)

func new4(s int32, m string) *Error { errno4xx++; return New(s, errno4xx, m) }
func new5(s int32, m string) *Error { errno5xx++; return New(s, errno5xx, m) }

// Server related
var (
	RequestParamError  = new4(400, "request parameter error") // 40000
	ServerUnknownError = new5(500, "server unknown error")    // 50000
)

// GitHub related
var (
	GithubQueryRateLimitError     = new5(500, "failed to query github rate limit")     // 50001
	GithubQueryApiResponseError   = new5(500, "failed to query github api response")   // 50002
	GithubQueryRepoIssuesError    = new5(500, "failed to query github repo issues")    // 50003
	GithubQueryIssueTimelineError = new5(500, "failed to query github issue timeline") // 50004
)

// SCUT related
var (
	ScutQueryJwNoticesError   = new5(500, "failed to query scut jw notices")   // 50003
	ScutQuerySeNoticesError   = new5(500, "failed to query scut se notices")   // 50004
	ScutQueryGrNoticesError   = new5(500, "failed to query scut gr notices")   // 50005
	ScutQueryGzicNoticesError = new5(500, "failed to query scut gzic notices") // 50006
)
