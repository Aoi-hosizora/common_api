package exception

import (
	"fmt"
	"github.com/Aoi-hosizora/ahlib-web/xgin"
	"github.com/Aoi-hosizora/ahlib-web/xgin/headers"
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
	Type   string `json:"type,omitempty"`
	Detail string `json:"detail,omitempty"`

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
		dto.RequestID = c.Writer.Header().Get(headers.XRequestID)
		dto.Request = xgin.DumpRequest(c, xgin.WithSecretHeaders(headers.Authorization))
	}
	return dto
}

func BuildFullErrorDto(err interface{}, c *gin.Context, skip uint) (*ErrorDto, xruntime.TraceStack) {
	dto := BuildBasicErrorDto(err, c) // include request info
	var stack xruntime.TraceStack
	stack, dto.Filename, dto.Funcname, dto.LineIndex, dto.Line = xruntime.RuntimeTraceStackWithInfo(skip + 1)
	dto.TraceStack = make([]string, 0, len(stack))
	for _, frame := range stack {
		dto.TraceStack = append(dto.TraceStack, frame.String())
	}
	return dto, stack
}
