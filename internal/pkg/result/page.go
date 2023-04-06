package result

import (
	"github.com/Aoi-hosizora/goapidoc"
)

func init() {
	goapidoc.AddDefinitions(
		goapidoc.NewDefinition("_Page", "page response").
			Generics("T").
			Properties(
				goapidoc.NewProperty("page", "integer", true, "current page"),
				goapidoc.NewProperty("limit", "integer", true, "page size"),
				goapidoc.NewProperty("total", "integer", true, "data count"),
				goapidoc.NewProperty("data", "T[]", true, "page data"),
			),
	)
}

type Page struct {
	Page  uint32 `json:"page"`
	Limit uint32 `json:"limit"`
	Total uint32 `json:"total"`
	Data  any    `json:"data"`
}

func NewPage(page, limit, total uint32, data any) *Page {
	return &Page{Page: page, Limit: limit, Total: total, Data: data}
}
