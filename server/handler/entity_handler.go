package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/refto/server/database/model"
	"github.com/refto/server/server/request"
	"github.com/refto/server/server/response"
	"github.com/refto/server/service/entity"
	"github.com/refto/server/service/topic"
)

func GetEntities(c *gin.Context) {
	var req request.FilterEntities
	if !bindQuery(c, &req) {
		return
	}

	var (
		definition *model.Entity
		err        error
	)
	if len(req.Topics) == 1 && req.Page < 2 {
		definition, err = entity.Definition(req.Topics[0])
		if err != nil {
			Abort(c, err)
			return
		}
	}

	data, count, err := entity.Filter(req)
	if err != nil {
		Abort(c, err)
		return
	}

	topics, err := topic.Common(req.Topics)
	if err != nil {
		Abort(c, err)
		return
	}

	resp := response.FilterEntities{
		Definition:    definition,
		Entities:      data,
		EntitiesCount: count,
		Topics:        topics,
	}

	c.JSON(http.StatusOK, resp)
}
