package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/refto/server/database/model"
	"github.com/refto/server/server/handler"
	"github.com/refto/server/server/request"
	"github.com/refto/server/service/repository"
)

func RequestRepository(idParam string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var elem model.Repository
		var err error
		if idParam == "path" {
			elem, err = repository.FindByPath(c.Param("path"))
		} else {
			var id int64
			if !handler.BindID(c, &id, idParam) {
				return
			}

			elem, err = repository.FindByID(id)
		}
		if err != nil {
			handler.Abort(c, err)
			return
		}

		request.SetRepository(c, elem)
		c.Next()
	}
}
