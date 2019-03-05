package customserver 

import (
	"time"
	"net/http"
	"tsstream/config"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func RunWebSocketServer() (err error){
	r := gin.Default()
	if nil == r {
		return 
	}

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	hub := NewWebSocketServer()
	go hub.RunServer()

	authorized := r.Group("/")
	authorized.Use(func(c *gin.Context) {	
	})
	{
		authorized.GET("/demo",func(c *gin.Context) {
			c.JSON(http.StatusOK,gin.H{"code": 200, "msg" :"Success","data": "Welcome dispatch system!"})
		})

		authorized.GET("/", func(c *gin.Context) {
			WebSocketConn(hub,c.Writer, c.Request,c)
		})
	}

	s := &http.Server{
		Addr:           ":" + config.CustomConf.Port,
		Handler:        r,
		ReadTimeout:    1 * time.Second,
		WriteTimeout:   1 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Println("Start HTTP Server In 0.0.0.0:",config.CustomConf.Port)

	s.ListenAndServe()
	return 
}

/**
* 运行客户tcp服务
*/
func RunTcpServer() (err error) {
	host := config.CustomConf.Bind
	port := config.CustomConf.Port

	tcpServer := GetTcpServer(host,port)
	tcpServer.RunCustomServer()

	return 
}