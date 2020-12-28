package dto

import (
	"github.com/Aoi-hosizora/common_api/src/model/vo"
	"github.com/Aoi-hosizora/goapidoc"
)

func init() {
	goapidoc.AddDefinitions(
		goapidoc.NewDefinition("ScutPostItemDto", "Scut post response").
			Properties(
				goapidoc.NewProperty("title", "string", true, "post title"),
				goapidoc.NewProperty("url", "string", true, "post url"),
				goapidoc.NewProperty("mobile_url", "string", true, "post url in mobile"),
				goapidoc.NewProperty("type", "string", true, "post type, is some specific strings"),
				goapidoc.NewProperty("date", "string#date", true, "post date, format: 0000-00-00"),
			),
	)
}

type ScutPostItemDto struct {
	Title     string `json:"title"`
	Url       string `json:"url"`
	MobileUrl string `json:"mobile_url"`
	Type      string `json:"type"`
	Date      string `json:"date"`
}

func BuildScutPostItemDto(item *vo.ScutPostItem) *ScutPostItemDto {
	return &ScutPostItemDto{
		Title:     item.Title,
		Url:       item.Url,
		MobileUrl: item.MobileUrl,
		Type:      item.Type,
		Date:      item.Date,
	}
}

func BuildScutPostItemDtos(items []*vo.ScutPostItem) []*ScutPostItemDto {
	out := make([]*ScutPostItemDto, len(items))
	for idx, item := range items {
		out[idx] = BuildScutPostItemDto(item)
	}
	return out
}
