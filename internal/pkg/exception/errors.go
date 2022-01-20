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
	GetGithubRateLimitError     = new5(500, "get github rate limit failed")
	GetGithubIssueTimelineError = new5(500, "get github issue timeline failed")
)

// SCUT related
var (
	GetScutJwError = new5(500, "get scut jw failed")
	GetScutSeError = new5(500, "get scut se failed")
)
