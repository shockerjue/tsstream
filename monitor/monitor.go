package monitor

import (
	"time"
	"strings"
	"io/ioutil"
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
		/**
		* 上传节点数据
		*/
		authorized.POST("/monitor",func(c *gin.Context) {
			PostMonitorFunc(c)
		})
		/**
		* 获取监控信息
		*/
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

func SendPostRequestHttp(url string, param string) (response string,err error) {
	postReq, err := http.NewRequest("POST",url,strings.NewReader(param))
    if err != nil {
        return 
    }

    postReq.Header.Set("Content-Type", "application/json;encoding=utf-8")

    client := &http.Client{}
    resp, err := client.Do(postReq)
    if err != nil {
        return 
	} 
	
	body, err := ioutil.ReadAll(resp.Body) 
	if err != nil {
		return "",err
	}

	defer resp.Body.Close()

	response = string(body)

	return 
}

var monitorChan chan MonitorInfo = make(chan MonitorInfo,1)
func RunMonitorTimer() {
	for {
		select {
		case monitor := <- monitorChan:
			data,err := monitor.Encode()
			if nil != err {
				return 
			}

			response, err := SendPostRequestHttp("",data)
			if nil == err {
				return 
			}

			log.Println(response)
		}
	}
}