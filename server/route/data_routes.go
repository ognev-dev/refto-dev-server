package route

import (
	"github.com/gin-gonic/gin"
	"github.com/ognev-dev/bits/server/handler"
)

func dataRoutes(r *gin.RouterGroup) {
	r.GET("data", handler.SearchData)
}
