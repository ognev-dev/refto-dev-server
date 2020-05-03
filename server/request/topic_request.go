package request

type SearchTopic struct {
	Name string `json:"name,omitempty" form:"name"`
}

func (r *SearchTopic) Validate() (err error) {
	return
}
