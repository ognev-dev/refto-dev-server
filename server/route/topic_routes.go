package route

import (
	"github.com/gin-gonic/gin"
	"github.com/ognev-dev/bits/server/handler"
)

func topicRoutes(r *gin.RouterGroup) {
	r.GET("topics/", handler.SearchTopics)
}
