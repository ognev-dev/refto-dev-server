package request

type FilterTopics struct {
	NoValidation
	Name string `json:"name,omitempty" form:"name"`
}
