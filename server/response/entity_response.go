package response

import "github.com/ognev-dev/bits/database/model"

type SearchEntity struct {
	Entities      []model.Entity `json:"entities"`
	EntitiesCount int            `json:"entities_count"`
	Topics        []model.Topic  `json:"topics"`
}
