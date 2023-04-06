package object

import (
	"sort"
)

type ScutNoticeItem struct {
	Title     string
	Url       string
	MobileUrl string // may equal to url
	Type      string
	Date      string // split by "-"
}

func NewScutNoticeItem(title, url, mobileUrl, typ, date string) *ScutNoticeItem {
	return &ScutNoticeItem{Title: title, Url: url, MobileUrl: mobileUrl, Type: typ, Date: date}
}

func SortScutNoticeItems(items []*ScutNoticeItem) {
	sort.Slice(items, func(i, j int) bool {
		return items[i].Date > items[j].Date // reverse
	})
}
