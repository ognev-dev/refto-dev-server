package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ognev-dev/bits/server/request"
	"github.com/ognev-dev/bits/server/response"
)

func SearchData(c *gin.Context) {
	var req request.SearchData
	if !bindQuery(c, &req) {
		return
	}

	//data, count, err := contact.Search(AuthUser(c).ID, req)
	//if err != nil {
	//	Abort(c, err)
	//	return
	//}

	resp := response.SearchData{
		Data:  data,
		Count: count,
	}

	c.JSON(http.StatusOK, resp)
}
