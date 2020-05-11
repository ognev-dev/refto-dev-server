package response

import "github.com/refto/server/database/model"

type SearchTopic struct {
	Data  []model.Topic `json:"data"`
	Count int           `json:"count"`
}
