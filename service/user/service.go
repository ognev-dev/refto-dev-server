package user

import (
	"errors"
	"time"

	"github.com/go-pg/pg/v9"
	"github.com/google/go-github/v32/github"
	"github.com/refto/server/database"
	"github.com/refto/server/database/model"
)

func ResolveFromGithubUser(gu *github.User, token string) (u model.User, err error) {
	if gu == nil {
		err = errors.New("github user is empty")
		return
	}

	if gu.GetID() == 0 {
		err = errors.New("github user's ID missing")
		return
	}

	cur, err := FindByGithubID(gu.GetID())
	if err != nil {
		return
	}

	// create new
	if cur == nil {
		u = model.User{
			Name:        gu.GetName(),
			AvatarURL:   gu.GetAvatarURL(),
			GithubID:    gu.GetID(),
			GithubToken: token,
			Email:       gu.GetEmail(),
			ActiveAt:    time.Now(),
		}
		err = Create(&u)
		return
	}

	// update existing
	cur.Name = gu.GetName()
	cur.AvatarURL = gu.GetAvatarURL()
	cur.GithubToken = token
	cur.Email = gu.GetEmail()
	cur.ActiveAt = time.Now()
	err = Update(cur)
	return
}

func FindByGithubID(id int64) (m *model.User, err error) {
	m = &model.User{}
	err = database.ORM().
		Model(m).
		Where("github_id = ?", id).
		First()

	if err == pg.ErrNoRows {
		m = nil
		err = nil
	}

	return
}

func Create(u *model.User) (err error) {
	err = database.ORM().Insert(u)
	return
}

func Update(u *model.User) (err error) {
	err = database.ORM().Update(u)
	return
}
