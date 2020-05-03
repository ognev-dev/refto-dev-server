package response

import "github.com/ognev-dev/bits/database/model"

type SearchTopic struct {
	Data  []model.Topic `json:"data"`
	Count int           `json:"count"`
}
