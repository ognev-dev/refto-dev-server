package topic

import (
	"github.com/ognev-dev/bits/database"
	"github.com/ognev-dev/bits/database/model"
	"github.com/ognev-dev/bits/server/request"
)

func Search(req request.SearchTopic) (data []model.Topic, count int, err error) {
	q := database.ORM().
		Model(&data)

	if req.Name != "" {
		q.Where("name ILIKE ?", req.Name+"%")
	}

	count, err = q.SelectAndCount()
	return
}
