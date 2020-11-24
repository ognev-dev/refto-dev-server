package route

import (
	"github.com/gin-gonic/gin"
	"github.com/refto/server/server/handler"
	"github.com/refto/server/server/middleware"
)

func collectionRoutes(r *gin.RouterGroup) {
	r.Group("collections/").
		GET("/", handler.GetCollections).
		POST("/", handler.CreateCollection)

	r.Group("collections/:id/").
		Use(middleware.RequestCollection("id")).
		PUT("/", handler.UpdateCollection).
		DELETE("/", handler.DeleteCollection)

	r.Group("collections/:token/").
		Use(middleware.RequestCollection("token")).
		GET("/", handler.GetCollectionByToken)

	r.Group("collections/:id/entities/:entity_id/").
		Use(middleware.RequestCollection("id")).
		Use(middleware.RequestEntity("entity_id")).
		POST("/", handler.AddEntityToCollection).
		DELETE("/", handler.RemoveEntityFromCollection)
}
