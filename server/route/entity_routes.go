package route

import (
	"github.com/gin-gonic/gin"
	"github.com/refto/server/server/handler"
)

func entityRoutes(r *gin.RouterGroup) {
	r.GET("entities/", handler.GetEntities)
}
