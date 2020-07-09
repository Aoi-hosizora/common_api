package exception

import (
	"fmt"
	"github.com/Aoi-hosizora/ahlib/xruntime"
	"github.com/gofiber/fiber"
	"time"
)

type Error struct {
	Code    int32
	Message string
}

func New(code int32, message string) *Error {
	return &Error{Code: code, Message: message}
}

func (e *Error) Error() string {
	return e.Message
}

type ErrorDto struct {
	Time    string `json:"time"`
	Type    string `json:"type"`
	Detail  string `json:"detail"`
	Request string `json:"request"`

	Filename string `json:"filename,omitempty"`
	Funcname string `json:"funcname,omitempty"`
	Line     int    `json:"line,omitempty"`
	Content  string `json:"content,omitempty"`
}

func BuildErrorDto(err interface{}, skip int, c *fiber.Ctx, print bool) *ErrorDto {
	now := time.Now().Format(time.RFC3339)
	errType := fmt.Sprintf("%T", err)
	errDetail := fmt.Sprintf("%v", err)
	if e, ok := err.(error); ok {
		errDetail = e.Error()
	}
	dto := &ErrorDto{
		Time:   now,
		Type:   errType,
		Detail: errDetail,
	}

	if c != nil {
		dto.Request = fmt.Sprintf("%s %s", c.Method(), c.OriginalURL())
	}

	if skip >= 0 {
		stacks, filename, funcname, line, content := xruntime.GetStackWithInfo(skip)
		dto.Filename = filename
		dto.Funcname = funcname
		dto.Line = line
		dto.Content = content
		if print {
			fmt.Println()
			xruntime.PrintStacks(stacks)
			fmt.Println()
		}
	}

	return dto
}
