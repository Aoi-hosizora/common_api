package param

import (
	"errors"
	"github.com/Aoi-hosizora/ahlib-mx/xgin"
	"github.com/Aoi-hosizora/ahlib/xnumber"
	"github.com/gin-gonic/gin"
)

type PageParam struct {
	Page  uint32
	Limit uint32
}

func BindQueryPage(c *gin.Context, defLimit, maxLimit uint32) *PageParam {
	page, err := xnumber.Atou32(c.Query("page"))
	if err != nil || page <= 0 {
		page = 1
	}

	limit, err := xnumber.Atou32(c.Query("limit"))
	if err != nil || limit <= 0 {
		limit = defLimit
	} else if limit > maxLimit {
		limit = maxLimit
	}

	return &PageParam{Page: page, Limit: limit}
}

func BindToken(c *gin.Context) string {
	token := c.GetHeader("Authorization")
	if token == "" {
		token = c.Query("token")
	}
	return token
}

func BindRouteID(c *gin.Context, name string) (uint64, error) {
	s := c.Param(name)
	id, err := xnumber.Atou64(s)
	if err != nil {
		return 0, xgin.NewRouterDecodeError(name, s, err, "")
	}
	if id <= 0 {
		err = errors.New("non-positive number")
		return 0, xgin.NewRouterDecodeError(name, s, err, "must be a positive number")
	}
	return id, nil
}

func BindBody[T any](c *gin.Context, obj T) (T, error) {
	err := c.ShouldBind(obj)
	if err != nil {
		var zero T
		return zero, err
	}
	return obj, nil
}
