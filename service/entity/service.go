package entity

import (
	"time"

	"github.com/go-pg/pg/v9"
	"github.com/ognev-dev/bits/database"
	"github.com/ognev-dev/bits/database/filter"
	"github.com/ognev-dev/bits/database/model"
	"github.com/ognev-dev/bits/server/request"
)

func Search(req request.SearchEntity) (data []model.Entity, count int, err error) {
	q := database.ORM().
		Model(&data).
		Apply(filter.PageFilter(req.Page, req.Limit))

	if len(req.Topics) > 0 {
		q.Join("JOIN entity_topics et ON et.entity_id=entity.id").
			Join("JOIN topics t ON et.topic_id=t.id").
			Group("entity.id")
	}

	// this will yield same results as with query from next comparison (topics > 1)
	// but this query just faster because no need for IN and HAVING
	// as it can be exact match
	if len(req.Topics) == 1 {
		q.Where("t.name = ?", req.Topics[0])
	}

	if len(req.Topics) > 1 {
		q.WhereIn("t.name IN (?)", req.Topics).
			Having("COUNT(t.id) = ?", len(req.Topics))
	}

	q.OrderExpr("updated_at DESC, created_at DESC")
	count, err = q.SelectAndCount()

	return
}

func CreateOrUpdate(elem *model.Entity) (err error) {
	old := model.Entity{}
	err = database.ORM().
		Model(&old).
		Where("token = ?", elem.Token).
		First()
	if err != nil && err != pg.ErrNoRows {
		return
	}

	// insert new
	if err == pg.ErrNoRows {
		err = database.ORM().Insert(elem)
		if err != nil {
			return
		}
	} else {
		old.Token = elem.Token
		old.Title = elem.Title
		old.Type = elem.Type
		old.DeletedAt = nil

		if old.Data != elem.Data {
			now := time.Now()
			old.Data = elem.Data
			old.UpdatedAt = &now
		}

		err = database.ORM().Update(&old)
		if err != nil {
			return
		}
		*elem = old
	}

	return nil
}
