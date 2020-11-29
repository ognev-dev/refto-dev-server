package factory

import (
	fake "github.com/brianvoe/gofakeit"
	"github.com/refto/server/database"
	"github.com/refto/server/database/model"
	"github.com/refto/server/service/collection"
)

func MakeCollection(opt ...model.Collection) (m model.Collection, err error) {
	if len(opt) == 1 {
		m = opt[0]
	}

	if m.Name == "" {
		m.Name = fake.Name()
	}
	if m.Token == "" {
		m.Token, err = collection.NewToken()
		if err != nil {
			return
		}
	}
	if m.UserID == 0 {
		if m.User == nil {
			userEl, err := CreateUser()
			if err != nil {
				return m, err
			}
			m.User = &userEl
		}
		m.UserID = m.User.ID
	}

	return
}

func CreateCollection(opt ...model.Collection) (m model.Collection, err error) {
	m, err = MakeCollection(opt...)
	if err != nil {
		return
	}

	err = database.ORM().Insert(&m)
	return
}
