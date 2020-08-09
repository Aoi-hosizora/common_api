package result

import (
	"github.com/Aoi-hosizora/ahlib-web/xdto"
	"github.com/Aoi-hosizora/ahlib-web/xgin"
	"github.com/Aoi-hosizora/common_api/src/common/exception"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type Result struct {
	Status  int32          `json:"-"`
	Code    int32          `json:"code"`
	Message string         `json:"message"`
	Data    interface{}    `json:"data,omitempty"`
	Error   *xdto.ErrorDto `json:"error,omitempty"`
}

func Status(status int32) *Result {
	message := http.StatusText(int(status))
	if status == 200 {
		message = "success"
	} else if message == "" {
		message = "unknown"
	}
	return &Result{
		Status:  status,
		Code:    status,
		Message: strings.ToLower(message),
	}
}

func Ok() *Result {
	return Status(http.StatusOK)
}

func Error(e *exception.Error) *Result {
	return Status(e.Code).SetMessage(e.Message)
}

func (r *Result) SetStatus(status int32) *Result {
	r.Status = status
	return r
}

func (r *Result) SetCode(code int32) *Result {
	r.Code = code
	return r
}

func (r *Result) SetMessage(message string) *Result {
	r.Message = strings.ToLower(message)
	return r
}

func (r *Result) SetData(data interface{}) *Result {
	r.Data = data
	return r
}

func (r *Result) SetPage(page int32, limit int32, total int32, data interface{}) *Result {
	r.Data = NewPage(page, limit, total, data)
	return r
}

func (r *Result) SetError(err error, c *gin.Context) *Result {
	if gin.Mode() == gin.DebugMode {
		r.Error = xgin.BuildBasicErrorDto(err, c)
	}
	return r
}

func (r *Result) JSON(c *gin.Context) {
	c.JSON(int(r.Status), r)
}

func (r *Result) XML(c *gin.Context) {
	c.XML(int(r.Status), r)
}
