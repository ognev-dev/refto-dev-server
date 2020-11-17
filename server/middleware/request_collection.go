package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/refto/server/server/handler"
	"github.com/refto/server/server/request"
	"github.com/refto/server/service/collection"
)

func RequestCollection(idParam string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var id int64
		if !handler.BindID(c, &id, idParam) {
			return
		}

		elem, err := collection.FindByID(id)
		if err != nil {
			handler.Abort(c, err)
			return
		}

		request.SetCollection(c, elem)
		c.Next()
	}
}
