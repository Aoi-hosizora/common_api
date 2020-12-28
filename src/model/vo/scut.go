package vo

import (
	"fmt"
)

type ScutPostItem struct {
	Title string
	Url   string
	Type  string
	Date  string
}

func (s *ScutPostItem) String() string {
	return fmt.Sprintf("%s %s %s %s", s.Title, s.Url, s.Type, s.Date)
}
