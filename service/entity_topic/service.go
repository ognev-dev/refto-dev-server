package entitytopic

import (
	"github.com/ognev-dev/bits/database"
	"github.com/ognev-dev/bits/database/model"
)

func Create(m model.EntityTopic) (err error) {
	err = database.ORM().Insert(&m)

	return
}
