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

	if len(req.Topics) > 0 {
		q.Join("LEFT JOIN data_topics dt ON dt.data_id=data.id").
			Join("LEFT JOIN topics t ON dt.topic_id=t.id")

		for _, v := range req.Topics {
			q.Where("t.name=?", v)
		}
	}

	q.Order("updated_at DESC, created_at DESC")

	count, err = q.SelectAndCount()

	return
}
