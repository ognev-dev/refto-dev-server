package request

type SearchEntity struct {
	Page   int      `json:"page,omitempty" form:"page"`
	Limit  int      `json:"limit,omitempty" form:"per_page"`
	Topics []string `json:"topics,omitempty" form:"topics"`
}

func (r *SearchEntity) Validate() (err error) {
	return
}
