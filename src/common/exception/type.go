package exception

type Error struct {
	Code    int32
	Message string
}

func (e *Error) Error() string {
	return e.Message
}

func New(code int32, message string) *Error {
	return &Error{Code: code, Message: message}
}
