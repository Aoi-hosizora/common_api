package result

type Page struct {
	Page  int32       `json:"page"`
	Limit int32       `json:"limit"`
	Total int32       `json:"total"`
	Data  interface{} `json:"data"`
}

func NewPage(page int32, limit int32, total int32, data interface{}) *Page {
	return &Page{Page: page, Limit: limit, Total: total, Data: data}
}
