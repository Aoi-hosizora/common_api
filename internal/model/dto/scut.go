package dto

import (
	"github.com/Aoi-hosizora/common_api/internal/model/biz"
	"github.com/Aoi-hosizora/goapidoc"
)

func init() {
	goapidoc.AddDefinitions(
		goapidoc.NewDefinition("ScutNoticeItemDto", "Scut notice item response").
			Properties(
				goapidoc.NewProperty("title", "string", true, "notice title"),
				goapidoc.NewProperty("url", "string", true, "notice url"),
				goapidoc.NewProperty("mobile_url", "string", true, "notice url in mobile"),
				goapidoc.NewProperty("type", "string", true, "notice type, is some specific strings"),
				goapidoc.NewProperty("date", "string#date", true, "notice date, format: 0000-00-00"),
			),
	)
}

type ScutNoticeItemDto struct {
	Title     string `json:"title"`
	Url       string `json:"url"`
	MobileUrl string `json:"mobile_url"`
	Type      string `json:"type"`
	Date      string `json:"date"`
}

func BuildScutNoticeItemDto(item *biz.ScutNoticeItem) *ScutNoticeItemDto {
	return &ScutNoticeItemDto{
		Title:     item.Title,
		Url:       item.Url,
		MobileUrl: item.MobileUrl,
		Type:      item.Type,
		Date:      item.Date,
	}
}

func BuildScutNoticeItemDtos(items []*biz.ScutNoticeItem) []*ScutNoticeItemDto {
	out := make([]*ScutNoticeItemDto, len(items))
	for idx, item := range items {
		out[idx] = BuildScutNoticeItemDto(item)
	}
	return out
}
