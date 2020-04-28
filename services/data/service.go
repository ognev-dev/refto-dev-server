package user

import (
	"github.com/ognev-dev/bits/cmd/server/errors"
	"github.com/ognev-dev/bits/database"
	"github.com/ognev-dev/bits/database/models"
)

func Create(elem *models.Data) (err error) {
	_, err = database.ORM().
		Model(elem).
		Insert()

	return
}

func FindByID(id string) (elem models.Data, err error) {
	err = database.ORM().
		Model(&elem).
		Where("id = ?", id).
		First()

	return
}

func Update(m *models.Data) (err error) {
	err = database.ORM().
		Update(m)

	return
}

func EmailVerified(id, email string) (err error) {
	if email != "" {
		var count int
		count, err = database.ORM().Model(&models.User{}).
			Where("email = ?", email).
			Where("id != ?", id).
			Count()
		if err != nil {
			return
		}
		if count > 0 {
			err = errors.EmailChangeEmailInUser
			return
		}
	}

	q := database.ORM().
		Model(&models.User{}).
		Set("email_verified = true").
		Where("id = ?", id)

	if email != "" {
		q.Set("email = ?", email)
	}

	_, err = q.Update()
	return
}
