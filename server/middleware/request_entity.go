package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/refto/server/server/handler"
	"github.com/refto/server/server/request"
	"github.com/refto/server/service/entity"
)

func RequestEntity(idParam string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var id int64
		if !handler.BindID(c, &id, idParam) {
			return
		}

		elem, err := entity.FindByID(id)
		if err != nil {
			handler.Abort(c, err)
			return
		}

		request.SetEntity(c, elem)
		c.Next()
	}
}
