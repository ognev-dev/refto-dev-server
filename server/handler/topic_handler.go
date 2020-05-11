package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/refto/server/server/request"
	"github.com/refto/server/server/response"
	"github.com/refto/server/service/topic"
)

func SearchTopics(c *gin.Context) {
	var req request.SearchTopic
	if !bindQuery(c, &req) {
		return
	}

	data, count, err := topic.Search(req)
	if err != nil {
		Abort(c, err)
		return
	}

	resp := response.SearchTopic{
		Data:  data,
		Count: count,
	}

	c.JSON(http.StatusOK, resp)
}
