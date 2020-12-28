package exception

var (
	cerr = int32(40000) - 1
	serr = int32(50000) - 1
)

func ce() int32 { cerr++; return cerr }
func se() int32 { serr++; return serr }

// request
var (
	RequestParamError   = New(400, ce(), "request param error")
	ServerRecoveryError = New(500, se(), "server unknown error")
)

// github
var (
	GetGithubRateLimitError     = New(500, se(), "get github rate limit failed")
	GetGithubIssueTimelineError = New(500, se(), "get github issue timeline failed")
	GetGithubRawPageError       = New(500, se(), "get github raw page error")
)
