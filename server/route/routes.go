package route

import (
	"io/ioutil"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/refto/server/config"
	"github.com/sirupsen/logrus"
)

func Register(r *gin.Engine) {
	conf := config.Get()

	r.Use(corsConfig())
	r.Use(static.Serve("/", static.LocalFile("./static", false)))

	api := r.Group(conf.Server.ApiBasePath)
	api.Use()

	apply(api,
		entityRoutes,
		topicRoutes,
	)

	r.NoRoute(func(c *gin.Context) {
		data, err := ioutil.ReadFile(filepath.Join("./static", static.INDEX))
		if err != nil {
			logrus.Error(err)
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		c.Data(http.StatusOK, "text/html", data)
	})
}

func apply(rg *gin.RouterGroup, routeFn ...func(*gin.RouterGroup)) {
	for _, fn := range routeFn {
		fn(rg)
	}
}

func corsConfig() gin.HandlerFunc {
	conf := cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		AllowCredentials: false,
		MaxAge:           24 * time.Hour,
	}

	return cors.New(conf)
}
