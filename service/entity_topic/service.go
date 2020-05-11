package entitytopic

import (
	"github.com/refto/server/database"
	"github.com/refto/server/database/model"
)

func Create(m model.EntityTopic) (err error) {
	err = database.ORM().Insert(&m)

	return
}
