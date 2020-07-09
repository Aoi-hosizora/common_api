package result

import (
	"encoding/json"
	"fmt"
	"github.com/Aoi-hosizora/common_api/src/common/exception"
	"github.com/Aoi-hosizora/common_api/src/config"
	"github.com/gofiber/fiber"
	"net/http"
	"strings"
)

type Result struct {
	Code    int32               `json:"code"`
	Message string              `json:"message"`
	Data    interface{}         `json:"data,omitempty"`
	Error   *exception.ErrorDto `json:"error,omitempty"`
}

type H map[string]interface{}

func (r *Result) Marshal() ([]byte, error) {
	if config.Configs.Meta.RunMode == "release" {
		r.Error = nil
	}
	return json.Marshal(r)
}

func (r *Result) String() string {
	b, err := r.Marshal()
	if err != nil {
		return fmt.Sprintf("%v", *r)
	}
	return string(b)
}

func Status(code int32) *Result {
	msg := http.StatusText(int(code))
	msg = strings.ToLower(msg)
	if code == 200 {
		msg = "success"
	}
	return &Result{Code: code, Message: msg}
}

func Ok() *Result {
	return Status(200)
}

func Error(err *exception.Error) *Result {
	return Status(err.Code).SetMessage(err.Message)
}

func (r *Result) SetCode(code int32) *Result {
	r.Code = code
	return r
}

func (r *Result) SetMessage(message string) *Result {
	r.Message = message
	return r
}

func (r *Result) SetData(data interface{}) *Result {
	r.Data = data
	return r
}

func (r *Result) SetError(err interface{}, c *fiber.Ctx) *Result {
	r.Error = exception.BuildErrorDto(err, -1, c, false)
	return r
}

func (r *Result) JSON(c *fiber.Ctx) {
	c.Fasthttp.Response.SetStatusCode(int(r.Code))
	c.Fasthttp.Response.Header.SetContentType(fiber.MIMEApplicationJSON)
	c.Fasthttp.Response.SetBodyString(r.String())
}
