package exception

var (
	errno4 = int32(40000) - 1
	errno5 = int32(50000) - 1
)

func new4(s int32, m string) *Error { errno4++; return New(s, errno4, m) }
func new5(s int32, m string) *Error { errno5++; return New(s, errno5, m) }

// Server related
var (
	RequestParamError  = new4(400, "request parameter error") // 40000
	ServerUnknownError = new5(500, "server unknown error")    // 50000
)

// GitHub related
var (
	GithubQueryRateLimitError     = new5(500, "failed to query github rate limit")     // 50001
	GithubQueryIssueTimelineError = new5(500, "failed to query github issue timeline") // 50002
)

// SCUT related
var (
	ScutQueryJwNoticesError   = new5(500, "failed to query scut jw notices")   // 50003
	ScutQuerySeNoticesError   = new5(500, "failed to query scut se notices")   // 50004
	ScutQueryGrNoticesError   = new5(500, "failed to query scut gr notices")   // 50005
	ScutQueryGzicNoticesError = new5(500, "failed to query scut gzic notices") // 50006
)
