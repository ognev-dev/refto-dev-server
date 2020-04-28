package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/ognev-dev/bits/server/handlers"
)

func dataRoutes(r *gin.RouterGroup) {
	r.GET("data", handlers.FindDataByTopics)
}
