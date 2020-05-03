package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ognev-dev/bits/server/request"
	"github.com/ognev-dev/bits/server/response"
	"github.com/ognev-dev/bits/service/topic"
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
