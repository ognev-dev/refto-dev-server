package topic

import (
	"github.com/go-pg/pg/v9"
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

func FirstOrCreate(name string) (elem model.Topic, err error) {
	err = database.ORM().
		Model(&elem).
		Where("name = ?", name).
		First()
	if err != nil && err != pg.ErrNoRows {
		return
	}

	if err == pg.ErrNoRows {
		elem.Name = name
		err = database.ORM().Insert(&elem)
	}

	return
}
