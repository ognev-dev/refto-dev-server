package handler

import (
	"net/http"

	"github.com/refto/server/database/model"

	"github.com/refto/server/service/repository"

	"github.com/gin-gonic/gin"
	"github.com/refto/server/server/request"
	"github.com/refto/server/server/response"
	"github.com/refto/server/service/collection"
)

func GetRepositories(c *gin.Context) {
	var req request.FilterRepositories
	if !bindQuery(c, &req) {
		return
	}

	req.Types = []model.RepositoryType{
		model.RepositoryTypeGlobal,
		model.RepositoryTypePublic,
	}
	data, count, err := repository.Filter(req)
	if err != nil {
		Abort(c, err)
		return
	}

	c.JSON(http.StatusOK, response.FilterRepositories{
		Data:  data,
		Count: count,
	})
}

func GetUserRepositories(c *gin.Context) {
	var req request.FilterRepositories
	if !bindQuery(c, &req) {
		return
	}

	req.UserID = request.User(c).ID
	data, count, err := repository.Filter(req)
	if err != nil {
		Abort(c, err)
		return
	}

	c.JSON(http.StatusOK, response.FilterRepositories{
		Data:  data,
		Count: count,
	})
}

func GetRepositoryByToken(c *gin.Context) {
	col, err := collection.FindByToken(c.Param("token"))
	if err != nil {
		Abort(c, err)
		return
	}

	c.JSON(http.StatusOK, col)
}

func CreateRepository(c *gin.Context) {
	var req request.CreateRepository
	if !bindJSON(c, &req) {
		return
	}

	elem := req.ToModel(c)
	secret, err := repository.Create(&elem)
	if err != nil {
		Abort(c, err)
		return
	}

	resp := response.CreateRepository{
		Secret: secret,
	}

	c.JSON(http.StatusCreated, resp)
}

func GetNewRepositorySecret(c *gin.Context) {
	userID := request.User(c).ID
	repo := request.Repository(c)
	if userID != repo.UserID {
		Abort(c, repository.ErrOwnershipViolation)
		return
	}

	secret, err := repository.NewSecret(repo.ID)
	if err != nil {
		Abort(c, err)
		return
	}

	resp := response.CreateRepository{
		Secret: secret,
	}

	c.JSON(http.StatusOK, resp)
}

func UpdateRepository(c *gin.Context) {
	//var req request.UpdateRepository
	//if !bindJSON(c, &req) {
	//	return
	//}
	//
	//elem := request.Repository(c)
	//req.ToModel(&elem)
	//err := collection.Update(&elem)
	//if err != nil {
	//	Abort(c, err)
	//	return
	//}

	//c.JSON(http.StatusOK, elem)
}

func DeleteRepository(c *gin.Context) {
	//if !validRequest(c, request.DeleteRepository{}) {
	//	return
	//}
	//
	//err := collection.Delete(request.Repository(c).ID)
	//if err != nil {
	//	Abort(c, err)
	//	return
	//}

	c.JSON(http.StatusOK, response.OK("Repository deleted"))
}
