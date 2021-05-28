package repository

import (
	"github.com/go-pg/pg/v9"
	"github.com/refto/server/database"
	"github.com/refto/server/database/model"
	"github.com/refto/server/errors"
	"github.com/refto/server/util"
)

var (
	ErrUserAlreadyAddedRepo = errors.Unprocessable("You already added this repository before")
	ErrRepoAlreadyClaimed   = errors.Unprocessable("Another user already added this repository")
	ErrRepoNotFoundByPath   = errors.NotFound("Repository is not found by path")
)

//func Filter(req request.FilterCollections) (data []model.Collection, count int, err error) {
//	q := database.ORM().
//		Model(&data).
//		Apply(filter.PageFilter(req.Page, req.Limit)).
//		Apply(filter.UserFilter(req.UserID)).
//		Order("created_at DESC")
//
//	if req.EntityID != 0 {
//		if req.Available {
//			q.Where("collection.id NOT IN (SELECT collection_id FROM collection_entities WHERE entity_id = ?)", req.EntityID)
//		} else {
//			q.Join("JOIN collection_entities ce ON ce.collection_id=collection.id").
//				Where("ce.entity_id = ?", req.EntityID)
//		}
//	}
//
//	if req.Name != "" {
//		q.Where("collection.name ILIKE ?", "%"+req.Name+"%")
//	}
//
//	if req.WithEntitiesCount {
//		q.Join("LEFT JOIN collection_entities ce2 ON ce2.collection_id=collection.id").
//			ColumnExpr("collection.*, COUNT(ce2.collection_id) AS entities_count").
//			Group("collection.id")
//
//	}
//
//	count, err = q.SelectAndCount()
//	return
//}
//

func Create(m *model.Repository) (secret string, err error) {
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

	secret = util.RandomString()
	m.Secret, err = util.HashPassword(secret)
	if err != nil {
		return
	}

	confirmed := false
	m.Confirmed = &confirmed

	err = database.ORM().Insert(m)
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

//
//func Update(elem *model.Collection) (err error) {
//	err = database.ORM().Update(elem)
//	return
//}
//
//func Delete(id int64) (err error) {
//	_, err = database.ORM().Model(&model.CollectionEntity{}).Where("collection_id = ?", id).Delete()
//	if err != nil {
//		return
//	}
//	_, err = database.ORM().Model(&model.Collection{}).Where("id = ?", id).Delete()
//
//	return
//}

func FindByID(id int64) (m model.Repository, err error) {
	err = database.ORM().
		Model(&m).
		Where("id = ?", id).
		First()

	return
}

//
//func FindByToken(token string) (m model.Collection, err error) {
//	err = database.ORM().
//		Model(&m).
//		Where("token = ?", token).
//		First()
//
//	if err == pg.ErrNoRows {
//		err = errors.New("collection not found")
//	}
//
//	return
//}
