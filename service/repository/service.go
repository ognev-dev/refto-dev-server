package repository

import (
	"strings"

	"github.com/go-pg/pg/v9"
	"github.com/refto/server/database"
	"github.com/refto/server/database/filter"
	"github.com/refto/server/database/model"
	"github.com/refto/server/errors"
	"github.com/refto/server/util"
)

var (
	ErrUserAlreadyAddedRepo = errors.Unprocessable("You already added this repository before")
	ErrRepoAlreadyClaimed   = errors.Unprocessable("Another user already added this repository")
	ErrRepoNotFoundByPath   = errors.NotFound("Repository is not found by path")
	ErrNotConfirmed         = errors.Unprocessable("Repository is not confirmed yet")
	ErrInvalidRepoPath      = func(given string) error {
		return errors.Unprocessable("Invalid repository path, must be in format of '{user}/{repo}', but '" + given + "' given")
	}
)

func Filter(req FilterParams) (data []model.Repository, count int, err error) {
	q := database.ORM().
		Model(&data).
		Apply(filter.PageFilter(req.Page, req.Limit)).
		Order("created_at DESC")

	if req.UserID > 0 {
		q.Apply(filter.UserFilter(req.UserID))
	}

	if req.Name != "" && len(req.Name) > 2 {
		q.Where("repository.name ILIKE ?", "%"+req.Name+"%")
	}
	if req.Path != "" && len(req.Name) > 2 {
		q.Where("repository.path ILIKE ?", "%"+req.Path+"%")
	}
	if len(req.Types) == 1 {
		q.Where("repository.type = ?", req.Types[0])
	}
	if len(req.Types) > 1 {
		q.WhereIn("repository.type IN (?)", req.Types)
	}
	if req.Confirmed != nil {
		q.Where("repository.confirmed IS ?", req.Confirmed)
	}

	count, err = q.SelectAndCount()
	return
}

func Create(m *model.Repository) (err error) {
	var existing model.Repository
	existing, err = FindByPath(m.Path)
	if err == nil {
		if existing.UserID == m.UserID {
			err = ErrUserAlreadyAddedRepo
		} else {
			err = ErrRepoAlreadyClaimed
		}
	}
	if err == ErrRepoNotFoundByPath {
		err = nil
	}
	if err != nil {
		return
	}

	m.Secret = util.RandomString()
	m.Confirmed = false

	err = database.ORM().Insert(m)
	return
}

func NewSecret(repoID int64) (secret string, err error) {
	secret = util.RandomString()
	err = UpdateSecret(repoID, secret)
	return
}

func UpdateSecret(repoID int64, secret string) (err error) {
	_, err = database.ORM().
		Model(&model.Repository{}).
		Where("id = ?", repoID).
		Set("secret = ?", secret).
		Update()

	return
}

func FindByPath(path string) (m model.Repository, err error) {
	err = database.ORM().
		Model(&m).
		Where("path = ?", path).
		First()

	if err == pg.ErrNoRows {
		err = ErrRepoNotFoundByPath
	}

	return
}

func Update(elem *model.Repository) (err error) {
	err = database.ORM().Update(elem)
	return
}

func Delete(id int64) (err error) {
	// TODO this can take some time if DB is loaded and busy
	// 	make it async or scheduled

	// TODO this can be unexpected for users who adds entities from this repo to their collections
	// 	it will be good to notify anyone related to this repo about this event

	_, err = database.ORM().
		Exec("DELETE FROM collection_entities WHERE entity_id IN (SELECT id FROM entities WHERE repo_id = ?)", id)
	if err != nil {
		return
	}
	_, err = database.ORM().
		Exec("DELETE FROM entity_topics WHERE entity_id IN (SELECT id FROM entities WHERE repo_id = ?)", id)
	if err != nil {
		return
	}

	_, err = database.ORM().
		Exec("DELETE FROM entities WHERE repo_id = ?", id)
	if err != nil {
		return
	}

	_, err = database.ORM().
		Exec("DELETE FROM topics WHERE repo_id = ?", id)
	if err != nil {
		return
	}

	_, err = database.ORM().
		Model(&model.Repository{}).
		Where("id = ?", id).
		Delete()

	return
}

func FindByID(id int64) (m model.Repository, err error) {
	err = database.ORM().
		Model(&m).
		Where("id = ?", id).
		First()

	return
}

// NormalizeRepoPath removes simple mistakes from potentially valid format
// (copying correct path from  GitHub's html page will result invalid path when pasting (whitespace around "/"))
// return error if path in invalid format
func NormalizeRepoPath(path string) (string, error) {
	path = strings.TrimPrefix(path, "/")
	path = strings.TrimSuffix(path, "/")
	parts := strings.Split(path, "/")
	if len(parts) != 2 {
		return "", ErrInvalidRepoPath(path)
	}

	path = strings.TrimSpace(parts[0]) + "/" + strings.TrimSpace(parts[1])
	return path, nil
}
