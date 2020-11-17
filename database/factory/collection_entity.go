package factory

import (
	"github.com/refto/server/database"
	"github.com/refto/server/database/model"
)

func MakeCollectionEntity(opt ...model.CollectionEntity) (m model.CollectionEntity, err error) {
	if len(opt) == 1 {
		m = opt[0]
	}

	if m.CollectionID == 0 {
		collectionElem, err := CreateCollection()
		if err != nil {
			return m, err
		}
		m.CollectionID = collectionElem.ID
	}
	if m.EntityID == 0 {
		entityElem, err := CreateEntity()
		if err != nil {
			return m, err
		}
		m.EntityID = entityElem.ID
	}

	return
}

func CreateCollectionEntity(opt ...model.CollectionEntity) (m model.CollectionEntity, err error) {
	m, err = MakeCollectionEntity(opt...)
	if err != nil {
		return
	}

	err = database.ORM().Insert(&m)
	return
}
