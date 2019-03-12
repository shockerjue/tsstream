package monitor

import (
	"time"
	"net/http"
	"tsstream/config"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

/**
* 接受HTTP流服务处理
*/

type MonitorServer struct {
	Host 	string
	Port 	string
}

func GetMonitorServer() *MonitorServer {
	monitor := &MonitorServer{Host: config.MonitorConf.Bind, Port: config.MonitorConf.Port}

	return monitor
}

func (this *MonitorServer)RunServer() (err error) {
	r := gin.Default()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())


	authorized := r.Group("/")
	authorized.Use(func(c *gin.Context) {	
	})
	{
		authorized.POST("/monitor",func(c *gin.Context) {
			PostMonitorFunc(c)
		})
		authorized.GET("/monitor",func(c *gin.Context){
			GetMonitorFunc(c)
		})
	}

	s := &http.Server{
		Addr:           ":" + this.Port,
		Handler:        r,
		ReadTimeout:    2 * time.Second, // 移除超时，这样可以连续不断的读取数据,//
		WriteTimeout:   0,
		MaxHeaderBytes: 1 << 20,
	}

	log.Println("Start monitor HTTP Server In ",this.Host," : ",this.Port)

	s.ListenAndServe()

	return 
}