package dto

import (
	"github.com/Aoi-hosizora/common_api/src/model/vo"
)

type ScutPostItemDto struct {
	Title string
	Url   string
	Type  string
	Date  string
}

func BuildScutPostItemDto(item *vo.ScutPostItem) *ScutPostItemDto {
	return &ScutPostItemDto{
		Title: item.Title,
		Url:   item.Url,
		Type:  item.Type,
		Date:  item.Date,
	}
}

func BuildScutPostItemDtos(items []*vo.ScutPostItem) []*ScutPostItemDto {
	out := make([]*ScutPostItemDto, len(items))
	for idx, item := range items {
		out[idx] = BuildScutPostItemDto(item)
	}
	return out
}
