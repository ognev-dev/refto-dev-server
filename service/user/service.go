package user

import (
	"errors"
	"time"

	"github.com/go-pg/pg/v9"
	"github.com/google/go-github/v32/github"
	"github.com/refto/server/database"
	"github.com/refto/server/database/model"
)

func ResolveFromGithubUser(gu *github.User, token string) (u *model.User, err error) {
	if gu == nil {
		err = errors.New("github user is empty")
		return
	}

	if gu.GetID() == 0 {
		err = errors.New("github user's ID missing")
		return
	}

	u, err = FindByGithubID(gu.GetID())
	if err != nil {
		return
	}

	// create new
	if u == nil {
		u = &model.User{
			Name:        gu.GetName(),
			Login:       gu.GetLogin(),
			AvatarURL:   gu.GetAvatarURL(),
			GithubID:    gu.GetID(),
			GithubToken: token,
			Email:       gu.GetEmail(),
			ActiveAt:    time.Now(),
		}
		err = Create(u)
		return
	}

	// update existing
	u.Name = gu.GetName()
	u.Login = gu.GetLogin()
	u.AvatarURL = gu.GetAvatarURL()
	u.GithubToken = token
	u.Email = gu.GetEmail()
	u.ActiveAt = time.Now()

	err = Update(u)

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

func SetActiveAt(userID int64) (err error) {
	_, err = database.ORM().
		Model(&model.User{}).
		Where("id = ?", userID).
		Set("active_at = ?", time.Now()).
		Update()

	return
}
