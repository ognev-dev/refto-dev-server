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

	u, err = FindByGithubID(gu.GetID())
	if err != nil && err != pg.ErrNoRows {
		return
	}

	// create new
	if err == pg.ErrNoRows {
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
	u.Name = gu.GetName()
	u.AvatarURL = gu.GetAvatarURL()
	u.GithubToken = token
	u.Email = gu.GetEmail()
	u.ActiveAt = time.Now()
	err = Update(&u)

	return
}

func FindByGithubID(id int64) (m model.User, err error) {
	err = database.ORM().
		Model(&m).
		Where("github_id = ?", id).
		First()

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
