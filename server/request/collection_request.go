package request

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/refto/server/database/model"
	se "github.com/refto/server/server/error"
)

const ctxCollectionKey = "collection_model"

func Collection(c *gin.Context) model.Collection {
	elem := c.MustGet(ctxCollectionKey)
	return elem.(model.Collection)
}

func SetCollection(c *gin.Context, m model.Collection) {
	c.Set(ctxCollectionKey, m)
}

type FilterCollections struct {
	NoValidation
	Pagination
	UserID int64 `json:"-" form:"-"`
}

type CreateCollection struct {
	Name    string `json:"name"`
	Private bool   `json:"private"`
}

func (r *CreateCollection) Validate(*gin.Context) (err error) {
	if strings.TrimSpace(r.Name) == "" {
		return se.Input{
			"name": "Name is required",
		}
	}

	return
}

func (r *CreateCollection) ToModel(c *gin.Context) (m model.Collection) {
	return model.Collection{
		Name:    r.Name,
		Private: r.Private,
		UserID:  User(c).ID,
	}
}

type UpdateCollection struct {
	Name    string `json:"name"`
	Private bool   `json:"private"`
}

func (r UpdateCollection) Validate(c *gin.Context) (err error) {
	if User(c).ID != Collection(c).UserID {
		err = errors.New("invalid user")
		return
	}

	if strings.TrimSpace(r.Name) == "" {
		return se.Input{
			"name": "Name is required",
		}
	}

	return
}

func (r UpdateCollection) ToModel(m *model.Collection) {
	m.Name = r.Name
	m.Private = r.Private
}

type DeleteCollection struct{}

func (r DeleteCollection) Validate(c *gin.Context) (err error) {
	if User(c).ID != Collection(c).UserID {
		err = errors.New("you can't simply delete collection that is not created by you")
		return
	}

	return
}

type AddEntityToCollection struct{}

func (r AddEntityToCollection) Validate(c *gin.Context) (err error) {
	if User(c).ID != Collection(c).UserID {
		err = errors.New("unable to add entity to collection: Ownership violation detected. ")
		return
	}

	return
}

type RemoveEntityFromCollection struct{}

func (r RemoveEntityFromCollection) Validate(c *gin.Context) (err error) {
	if User(c).ID != Collection(c).UserID {
		err = errors.New("unable to remove entity from collection: Owner of collection will not be happy")
		return
	}

	return
}
