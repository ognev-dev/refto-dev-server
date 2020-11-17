package factory

import (
	fake "github.com/brianvoe/gofakeit"
	"github.com/refto/server/database"
	"github.com/refto/server/database/model"
)

func MakeTopic(opt ...model.Topic) (m model.Topic, err error) {
	if len(opt) == 1 {
		m = opt[0]
	}

	if m.Name == "" {
		m.Name = fake.Name()
	}

	return
}

func CreateTopic(opt ...model.Topic) (m model.Topic, err error) {
	m, err = MakeTopic(opt...)
	if err != nil {
		return
	}

	err = database.ORM().Insert(&m)
	return
}
