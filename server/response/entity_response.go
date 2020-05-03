package response

import "github.com/ognev-dev/bits/database/model"

type SearchEntity struct {
	Data  []model.Entity `json:"data"`
	Count int            `json:"count"`
}
