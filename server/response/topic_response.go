package response

import "github.com/refto/server/database/model"

type FilterTopics struct {
	Data  []model.Topic `json:"data"`
	Count int           `json:"count"`
}
