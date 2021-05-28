package test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/refto/server/service/repository"

	"github.com/refto/server/database/factory"

	"github.com/refto/server/server/response"

	"github.com/brianvoe/gofakeit"
	"github.com/refto/server/database/model"
	"github.com/refto/server/server/request"
	. "github.com/refto/server/test/apitest"
	"github.com/refto/server/test/assert"
	"github.com/refto/server/util"
)

func TestCreateRepository(t *testing.T) {
	Authorise(t)

	req := request.CreateRepository{
		Path:        util.RandomString() + "/" + util.RandomString(),
		Name:        gofakeit.Name(),
		Description: gofakeit.Name(),
		Type:        model.RepositoryTypeGlobal,
	}

	var resp response.CreateRepository
	TestCreate(t, "repositories", req, &resp)

	assert.True(t, resp.Secret != "")
	assert.DatabaseHas(t, "repositories", util.M{
		"path":        req.Path,
		"user_id":     AuthUser.ID,
		"name":        req.Name,
		"description": req.Description,
		"type":        req.Type,
		"confirmed":   false,
	})
}

// Since repo path is unique, adding repo with existing path should raise error
// TODO add repo transfer to another user using secret
func TestCreateRepository_Existing(t *testing.T) {
	Authorise(t)

	m, err := factory.CreateRepository()
	assert.NotError(t, err)

	req := request.CreateRepository{
		Path:        m.Path,
		Name:        gofakeit.Name(),
		Description: gofakeit.Name(),
		Type:        model.RepositoryTypeGlobal,
	}
	resp, _ := TestCreate422(t, "repositories", req)
	assert.Equals(t, resp.Error, repository.ErrRepoAlreadyClaimed.Error())
}

func TestRepositoryGetNewSecret(t *testing.T) {
	Authorise(t)

	m, err := factory.CreateRepository(model.Repository{
		UserID: AuthUser.ID,
	})
	assert.NotError(t, err)

	var resp response.CreateRepository
	POST(t, Request{
		Path:         fmt.Sprintf("repositories/%d/secret/", m.ID),
		Body:         nil,
		BindResponse: &resp,
		AssertStatus: http.StatusOK,
	})
	assert.True(t, resp.Secret != "")
}

func TestGetUserRepositories(t *testing.T) {
	Authorise(t)

	m1, err := factory.CreateRepository(model.Repository{UserID: AuthUser.ID})
	assert.NotError(t, err)
	m2, err := factory.CreateRepository(model.Repository{UserID: AuthUser.ID})
	assert.NotError(t, err)
	_, err = factory.CreateRepository() // not user's
	assert.NotError(t, err)

	var req request.FilterRepositories
	var resp response.FilterRepositories
	TestFilter(t, "user/repositories/", req, &resp)
	assert.Equals(t, 2, resp.Count)

	for _, el := range resp.Data {
		if el.ID == m1.ID {
			continue
		}
		if el.ID == m2.ID {
			continue
		}

		t.Fatalf("invalid element in response: %v", el)
	}
}
