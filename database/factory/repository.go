package factory

import (
	fake "github.com/brianvoe/gofakeit"
	"github.com/refto/server/database"
	"github.com/refto/server/database/model"
	"github.com/refto/server/util"
)

func MakeRepository(opt ...model.Repository) (m model.Repository, err error) {
	if len(opt) == 1 {
		m = opt[0]
	}

	if m.Name == "" {
		m.Name = fake.Name()
	}
	if m.Description == "" {
		m.Description = fake.Name()
	}
	if m.User != nil {
		m.UserID = m.User.ID
	}
	if m.UserID == 0 {
		var u model.User
		u, err = CreateUser()
		if err != nil {
			return
		}
		m.UserID = u.ID
		m.User = &u
	}
	if m.Path == "" {
		m.Path = util.RandomString() + "/" + util.RandomString()
	}
	if m.Type == "" {
		m.Type = model.RepositoryType(fake.RandString(model.RepositoryTypesList))
	}
	if m.Confirmed == nil {
		confirmed := util.RandomBool()
		m.Confirmed = &confirmed
	}
	if m.Secret == "" {
		m.Secret, err = util.HashPassword(util.RandomString(10))
		if err != nil {
			return
		}
	}

	return
}

func CreateRepository(opt ...model.Repository) (m model.Repository, err error) {
	m, err = MakeRepository(opt...)
	if err != nil {
		return
	}

	err = database.ORM().Insert(&m)
	return
}
