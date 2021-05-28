package route

import (
	"github.com/gin-gonic/gin"
	"github.com/refto/server/server/handler"
	"github.com/refto/server/server/middleware"
)

func repositoryRoutes(r *gin.RouterGroup) {
	r.Group("repositories/").
		GET("/", handler.GetRepositories).
		POST("/", handler.CreateRepository)

	r.Group("repositories/:id/").
		Use(middleware.RequestRepository("id")).
		PUT("/", handler.UpdateCollection).
		DELETE("/", handler.DeleteCollection)

	r.Group("repositories/:token/").
		Use(middleware.RequestRepository("token")).
		GET("/", handler.GetCollectionByToken)

}
