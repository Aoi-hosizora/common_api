package exception

// request
var (
	RequestParamError   = New(400, "request param error")
	ServerRecoveryError = New(500, "server unknown error")
	GetGithubError      = New(500, "failed to get github response")
)
