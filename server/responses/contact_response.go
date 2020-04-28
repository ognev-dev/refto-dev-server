package responses

import "github.com/ognev-dev/bits/database/models"

type SearchData struct {
	Data  []models.Data `json:"data"`
	Count int           `json:"count"`
}
