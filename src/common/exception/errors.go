package exception

// request
var (
	ServerRecoveryError = New(500, "server unknown error")
)

// response
var (
	GetGithubError = New(500, "failed to get github response")
)
