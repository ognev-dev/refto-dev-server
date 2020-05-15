package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/gin-gonic/autotls"
	"github.com/gin-gonic/gin"
	"github.com/refto/server/config"
	"github.com/refto/server/logger"
	"github.com/refto/server/server/route"
)

func main() {
	conf := config.Get()
	logger.Setup()

	log.Println("refto." + conf.AppEnv + " serving at " + conf.Server.Host + ":" + conf.Server.Port)

	if config.IsReleaseEnv() {
		gin.SetMode(conf.AppEnv)
	}

	r := gin.Default()
	route.Register(r)

	if !config.IsReleaseEnv() {
		server := &http.Server{
			Addr:    conf.Server.Host + ":" + conf.Server.Port,
			Handler: r,
		}

		quit := make(chan os.Signal)
		signal.Notify(quit, os.Interrupt)

		go func() {
			<-quit
			if err := server.Close(); err != nil {
				log.Fatal(err.Error())
			}
		}()

		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal(err.Error())
		}
	} else {
		//On Linux, you can use setcap to grant your binary the permission to bind low ports:
		//$ sudo setcap cap_net_bind_service=+ep /path/to/your/binary
		err := autotls.Run(r, conf.Server.Host)
		if err != nil && err != http.ErrServerClosed {
			log.Fatal(err.Error())
		}
	}

	log.Println("refto." + conf.AppEnv + " server closed")
}
