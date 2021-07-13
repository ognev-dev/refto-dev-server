package filter

import (
	"github.com/go-pg/pg/v9/orm"
	_ "github.com/lib/pq"
)

const DefaultLimit = 25

type Fn func(*orm.Query) (*orm.Query, error)

func Apply(q *orm.Query, filters ...Fn) {
	for _, fn := range filters {
		q.Apply(fn)
	}
}

func PageFilter(page, limit int) Fn {
	if limit == 0 {
		limit = DefaultLimit
	}
	offset := 0
	if page > 1 {
		offset = (page - 1) * limit
	}

	return func(q *orm.Query) (*orm.Query, error) {
		if limit > 0 {
			q.Limit(limit)
		}
		if offset > 0 {
			q.Offset(offset)
		}
		return q, nil
	}
}

func UserFilter(userID int64) Fn {
	return func(q *orm.Query) (*orm.Query, error) {
		q.Where("user_id = ?", userID)

		return q, nil
	}
}

func TrashedFilter(trashed bool, optTable ...string) Fn {
	col := "deleted_at"
	if len(optTable) == 1 {
		col = optTable[0] + "." + col
	}
	where := col + " is"
	if trashed {
		where += " not"
	}
	where += " null"

	return func(q *orm.Query) (*orm.Query, error) {
		q.Where(where)

		return q, nil
	}
}
