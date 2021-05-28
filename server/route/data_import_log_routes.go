package route

import (
	"github.com/gin-gonic/gin"
	"github.com/refto/server/server/handler"
)

func dataImportLogRoutes(r *gin.RouterGroup) {
	r.GET("data-import-log/", handler.GetTopics)
}
