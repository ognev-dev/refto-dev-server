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
		if m.Collection == nil {
			m.Collection = &model.Collection{}
			*m.Collection, err = CreateCollection()
			if err != nil {
				return m, err
			}
		}
		m.CollectionID = m.Collection.ID
	}
	if m.EntityID == 0 {
		if m.Entity == nil {
			m.Entity = &model.Entity{}
			*m.Entity, err = CreateEntity()
			if err != nil {
				return m, err
			}
		}

		m.EntityID = m.Entity.ID
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
