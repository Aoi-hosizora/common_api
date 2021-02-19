package exception

import (
	"fmt"
	"github.com/Aoi-hosizora/ahlib-web/xgin"
	"github.com/Aoi-hosizora/ahlib/xcolor"
	"github.com/Aoi-hosizora/ahlib/xruntime"
	"github.com/gin-gonic/gin"
	"time"
)

type Error struct {
	Status  int32
	Code    int32
	Message string
}

func (e *Error) Error() string {
	return e.Message
}

func New(status int32, code int32, message string) *Error {
	return &Error{Status: status, Code: code, Message: message}
}

type ErrorDto struct {
	Time   string `json:"time"`
	Type   string `json:"type"`
	Detail string `json:"detail"`

	RequestID string   `json:"request_id,omitempty"`
	Request   []string `json:"request,omitempty"`

	Filename   string   `json:"filename,omitempty"`
	Funcname   string   `json:"funcname,omitempty"`
	LineIndex  int      `json:"line_index,omitempty"`
	Line       string   `json:"line,omitempty"`
	TraceStack []string `json:"trace_stack,omitempty"`
}

func BuildBasicErrorDto(err interface{}, c *gin.Context) *ErrorDto {
	now := time.Now().Format(time.RFC3339)
	typ := fmt.Sprintf("%T", err)
	detail := fmt.Sprintf("%v", err)
	dto := &ErrorDto{Time: now, Type: typ, Detail: detail}
	if c != nil {
		dto.RequestID = c.Writer.Header().Get("X-Request-ID")
		dto.Request = xgin.DumpRequest(c, xgin.WithSecretHeaders("Authorization"))
	}
	return dto
}

func BuildFullErrorDto(err interface{}, c *gin.Context) *ErrorDto {
	skip := 3
	dto := BuildBasicErrorDto(err, c)

	var stack xruntime.TraceStack
	stack, dto.Filename, dto.Funcname, dto.LineIndex, dto.Line = xruntime.RuntimeTraceStackWithInfo(skip)
	dto.TraceStack = make([]string, len(stack))
	for idx, frame := range stack {
		dto.TraceStack[idx] = frame.String()
	}

	fmt.Println()
	fmt.Println(xcolor.BrightRed.Sprintf(stack.String()))
	fmt.Println()
	return dto
}
