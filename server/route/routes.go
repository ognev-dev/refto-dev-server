package route

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/ognev-dev/bits/config"
)

func Register(r *gin.Engine) {
	conf := config.Get()

	r.Use(corsConfig())

	api := r.Group(conf.Server.ApiBasePath)
	api.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "bits."+conf.AppEnv)
	})

	apply(api,
		dataRoutes,
	)
}

func apply(rg *gin.RouterGroup, routesFn ...func(*gin.RouterGroup)) {
	for _, fn := range routesFn {
		fn(rg)
	}
}

func corsConfig() gin.HandlerFunc {
	conf := cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "X-Client", "Authorization"},
		AllowCredentials: false,
		MaxAge:           24 * time.Hour,
	}

	return cors.New(conf)
}
