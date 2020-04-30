package request

type SearchData struct {
	Page   int      `json:"page" form:"page"`
	Limit  int      `json:"limit" form:"per_page"`
	Topics []string `json:"query" form:"query"`
}

func (r *SearchData) Validate() (err error) {
	return
}
