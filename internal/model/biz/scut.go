package biz

import (
	"sort"
)

type ScutNoticeItem struct {
	Title     string
	Url       string
	MobileUrl string
	Type      string
	Date      string
}

func NewScutNoticeItem(title, url, mobileUrl, typ, date string) *ScutNoticeItem {
	return &ScutNoticeItem{Title: title, Url: url, MobileUrl: mobileUrl, Type: typ, Date: date}
}

func SortScutNoticeItems(items []*ScutNoticeItem) {
	sort.Slice(items, func(i, j int) bool {
		return items[i].Date > items[j].Date // reverse
	})
}
