
package router

import (
	"tsstream/controller"
	"github.com/gin-gonic/gin"
	// log "github.com/sirupsen/logrus"
)

func Router() *gin.Engine {
	r := gin.Default()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// hub := chat.NewHub()
	// go hub.Run()

	/**
	* router
	*/
	authorized := r.Group("/")
	authorized.Use(func(c *gin.Context) {	
	})
	{
		/*
		* login
		*/
		authorized.POST("/stream/:key",func(c *gin.Context) {
			controller.OnRecvStream(c)
		})

		/**
		* websocket connect
		*/
		authorized.GET("/wellchat", func(c *gin.Context) {
			// chat.ServeWs(hub,c.Writer, c.Request,c)
		})
	}

	return r	
}
