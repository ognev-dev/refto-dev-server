package route

import (
	"github.com/gin-gonic/gin"
	"github.com/ognev-dev/bits/server/handler"
)

func entityRoutes(r *gin.RouterGroup) {
	r.GET("entities/", handler.SearchEntities)
}
