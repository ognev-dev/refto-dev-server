package response

import "github.com/refto/server/database/model"

type FilterEntities struct {
	Definition    *model.Entity  `json:"definition"`
	Entities      []model.Entity `json:"entities"`
	EntitiesCount int            `json:"entities_count"`
	Topics        []string       `json:"topics"`
}
