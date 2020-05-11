package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/refto/server/server/request"
	"github.com/refto/server/server/response"
	"github.com/refto/server/service/entity"
	"github.com/refto/server/service/topic"
)

func SearchEntities(c *gin.Context) {
	var req request.SearchEntity
	if !bindQuery(c, &req) {
		return
	}

	data, count, err := entity.Search(req)
	if err != nil {
		Abort(c, err)
		return
	}

	topics, err := topic.Common(req.Topics)
	if err != nil {
		Abort(c, err)
		return
	}

	resp := response.SearchEntity{
		Entities:      data,
		EntitiesCount: count,
		Topics:        topics,
	}

	c.JSON(http.StatusOK, resp)
}
