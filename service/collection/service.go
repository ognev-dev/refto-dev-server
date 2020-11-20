package collection

import (
	"errors"
	"math/rand"

	"github.com/go-pg/pg/v9"
	"github.com/refto/server/database"
	"github.com/refto/server/database/filter"
	"github.com/refto/server/database/model"
	"github.com/refto/server/server/request"
)

func Filter(req request.FilterCollections) (data []model.Collection, count int, err error) {
	q := database.ORM().
		Model(&data).
		Apply(filter.PageFilter(req.Page, req.Limit)).
		Apply(filter.UserFilter(req.UserID)).
		Order("created_at DESC")

	if req.EntityID != 0 {
		q.Join("JOIN collection_entities ce ON ce.collection_id=collection.id").
			Where("ce.entity_id = ?", req.EntityID)
	}

	count, err = q.SelectAndCount()
	return
}

func Create(elem *model.Collection) (err error) {
	elem.Token, err = NewToken()
	if err != nil {
		return
	}

	err = database.ORM().Insert(elem)
	return
}

func Update(elem *model.Collection) (err error) {
	err = database.ORM().Update(elem)
	return
}

func Delete(id int64) (err error) {
	_, err = database.ORM().Model(&model.CollectionEntity{}).Where("collection_id = ?", id).Delete()
	if err != nil {
		return
	}
	_, err = database.ORM().Model(&model.Collection{}).Where("id = ?", id).Delete()

	return
}

func FindByID(id int64) (m model.Collection, err error) {
	err = database.ORM().
		Model(&m).
		Where("id = ?", id).
		First()

	if err == pg.ErrNoRows {
		err = errors.New("collection not found")
	}

	return
}

func AddEntity(collectionID, entityID int64) (err error) {
	err = database.ORM().Insert(&model.CollectionEntity{
		CollectionID: collectionID,
		EntityID:     entityID,
	})

	return
}

func RemoveEntity(collectionID, entityID int64) (err error) {
	_, err = database.ORM().
		Model(&model.CollectionEntity{}).
		Where("collection_id = ?", collectionID).
		Where("entity_id = ?", entityID).
		Delete()

	return
}

// Creating unique token that is only 1 lowercase char
// and increase it's length on each failed attempt
func NewToken() (token string, err error) {
	data := []byte("abcdefghijklmnopqrstuvwxyz")
	for length := 1; ; length++ {
		t := make([]byte, length)
		for char := 0; char < length; char++ {
			t[char] = data[rand.Intn(len(data))]
		}
		token = string(t)
		err = database.ORM().
			Model(&model.Collection{}).
			Where("token = ?", token).
			First()

		if err == pg.ErrNoRows {
			return token, nil
		}

		if err != nil {
			return "", err
		}
	}
}
