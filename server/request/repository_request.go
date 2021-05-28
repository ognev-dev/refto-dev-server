package request

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/refto/server/database/model"
	"github.com/refto/server/errors"
	"github.com/refto/server/util"
)

var (
	errRepoUpdateWrongUser = errors.Unprocessable("you can't simply update repository that is not created by you")
	errRepoDeleteWrongUser = errors.Unprocessable("you can't simply delete repository that is not created by you")
)

const ctxRepositoryKey = "repository_model"

func Repository(c *gin.Context) model.Repository {
	elem := c.MustGet(ctxRepositoryKey)
	return elem.(model.Repository)
}

func SetRepository(c *gin.Context, m model.Repository) {
	c.Set(ctxRepositoryKey, m)
}

type FilterRepositories struct {
	NoValidation
	Pagination
	Path string `json:"path" form:"path"`
	Name string `json:"name" form:"name"`
}

type CreateRepository struct {
	Path        string               `json:"path"`
	Name        string               `json:"name"`
	Description string               `json:"description"`
	Type        model.RepositoryType `json:"type"`
}

func (r *CreateRepository) Validate(*gin.Context) (err error) {
	errs := errors.NewInput()
	errs.AddIf(util.IsEmptyString(r.Name), "name", "Name is required")
	errs.AddIf(util.IsEmptyString(r.Description), "description", "Description is required")
	errs.AddIf(util.IsEmptyString(r.Path), "path", "Path is required")
	errs.AddIf(util.IsEmptyString(r.Type.String()), "type", "Type is required")
	errs.AddIf(!r.Type.IsValid(), "type", fmt.Sprintf("Invalid type: '%s'", r.Type))

	if errs.Has() {
		return errs
	}

	return nil
}

func (r *CreateRepository) ToModel(c *gin.Context) (m model.Repository) {
	return model.Repository{
		UserID:      User(c).ID,
		Path:        r.Path,
		Name:        r.Name,
		Description: r.Description,
		Type:        r.Type,
	}
}

type UpdateRepository struct {
	Name        string               `json:"name"`
	Description string               `json:"description"`
	Type        model.RepositoryType `json:"type"`
}

func (r UpdateRepository) Validate(c *gin.Context) (err error) {
	if User(c).ID != Repository(c).UserID {
		return errRepoUpdateWrongUser
	}

	errs := errors.NewInput()

	if r.Type != "" && !r.Type.IsValid() {
		errs["type"] = fmt.Sprintf("invalid type: '%s'", r.Type)
	}

	if len(errs) > 0 {
		return errs
	}

	return nil
}

func (r UpdateRepository) ToModel(m *model.Repository) {
	if r.Name != "" {
		m.Name = r.Name
	}
	if r.Description != "" {
		m.Name = r.Description
	}
	if r.Type != "" {
		m.Type = r.Type
	}
}

type DeleteRepository struct{}

func (r DeleteRepository) Validate(c *gin.Context) error {
	if User(c).ID != Repository(c).UserID {
		return errRepoDeleteWrongUser
	}

	return nil
}
