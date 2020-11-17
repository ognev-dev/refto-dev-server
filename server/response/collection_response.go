package response

import "github.com/refto/server/database/model"

type FilterCollections struct {
	Data  []model.Collection `json:"data"`
	Count int                `json:"count"`
}
