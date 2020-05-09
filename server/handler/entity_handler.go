package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ognev-dev/bits/server/request"
	"github.com/ognev-dev/bits/server/response"
	"github.com/ognev-dev/bits/service/entity"
	"github.com/ognev-dev/bits/service/topic"
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
