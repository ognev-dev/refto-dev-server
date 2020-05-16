package route

import (
	"github.com/gin-gonic/gin"
	"github.com/refto/server/server/handler"
)

func webHookRoutes(r *gin.RouterGroup) {
	r.POST("hooks/data-pushed/", handler.ImportDataFromRepoByGitHubWebHook)
}
