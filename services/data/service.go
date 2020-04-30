package data

import (
	"github.com/ognev-dev/bits/database"
	"github.com/ognev-dev/bits/database/filters"
	"github.com/ognev-dev/bits/database/models"
	"github.com/ognev-dev/bits/server/request"
)

func Search(req request.SearchData) (data []models.Data, count int, err error) {
	q := database.ORM().
		Model(&data).
		Apply(filters.PageFilter(req.Page, req.Limit))

	// todo filtration by topic

	count, err = q.SelectAndCount()

	return
}
