package route

import (
	"github.com/gin-gonic/gin"
	"github.com/refto/server/server/handler"
)

func userRoutes(r *gin.RouterGroup) {
	r.POST("user/login/", handler.LoginWithGithub)
	r.GET("user/repositories/", handler.GetUserRepositories)
}
