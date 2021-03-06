package api

import (
	"github.com/swaggo/swag"
	"io/ioutil"
)

type swagger struct{}

func (s *swagger) ReadDoc() string {
	f, err := ioutil.ReadFile("./api/spec.json")
	if err != nil {
		return ""
	}
	return string(f)
}

func RegisterSwag() {
	swag.Register(swag.Name, &swagger{})
}
