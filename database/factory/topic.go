package factory

import (
	fake "github.com/brianvoe/gofakeit"
	"github.com/ognev-dev/bits/database"
	"github.com/ognev-dev/bits/database/model"
)

func MakeTopic(mOpt ...model.Topic) (m model.Topic, err error) {
	if len(mOpt) == 1 {
		m = mOpt[0]
	}

	if m.Name == "" {
		m.Name = fake.Name()
	}

	return
}

func CreateTopic(mOpt ...model.Topic) (m model.Topic, err error) {
	m, err = MakeTopic(mOpt...)
	if err != nil {
		return
	}

	err = database.ORM().Insert(&m)
	return
}
