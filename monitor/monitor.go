package monitor

import (
	"fmt"
	"time"
	"strings"
	"io/ioutil"
	"net/http"
	"tsstream/config"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

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

func UploadMonitor(packages,connects int) {
	if 1 < len(monitorChan) {
		return 
	}

	if time.Now().Unix() < (currentTime + 60) {
		return
	}
	currentTime = time.Now().Unix()

	genesis := false
	if "normal" == config.AppConf.RunMode {
		genesis = true
	}

	bind := config.BackStreamConf.Bind
	port := config.BackStreamConf.Port
	if "extra" == config.AppConf.RunMode {
		bind = config.PushStreamConf.Bind
		port = config.PushStreamConf.Port
	} 

	var monitor MonitorInfo
	monitor.NodeInfo.Name = config.AppConf.RunMode
	monitor.NodeInfo.Connects = connects
	monitor.NodeInfo.Bind = bind
	monitor.NodeInfo.Port = port
	monitor.NodeInfo.Hash = ""

	monitor.Packages = packages
	monitor.Genesis = genesis

	monitor.NextNode = make([]Node,0)
	protocol := config.DispatchConf.Protocol
	if "UDP" == protocol || "TCP" == protocol {
		hosts := strings.Split(config.DispatchConf.Hosts,",")
		ports := strings.Split(config.DispatchConf.Ports,",")

		if len(ports) != len(hosts) {
			return 
		}

		for k,_ := range hosts {
			var node Node 
			node.Name = fmt.Sprintf("extra%d", k + 1)
			node.Connects = 0
			node.Bind = hosts[k]
			node.Port = ports[k]
			node.Hash = ""

			monitor.NextNode = append(monitor.NextNode, node)
		}
	}
	
	if "CUSTOM" == protocol {
		var node Node 
		node.Name = "WEBSOCKET"
		node.Connects = 0
		node.Bind = config.CustomConf.Bind
		node.Port = config.CustomConf.Port
		node.Hash = ""

		monitor.NextNode = append(monitor.NextNode, node)
	}

	if 1 < len(monitorChan) {
		return 
	}

	monitorChan <- monitor

	return 
}

var currentTime int64 = 0
var monitorChan chan MonitorInfo = make(chan MonitorInfo,2)
func RunMonitorTimer() {
	for {
		select {
		case monitor := <- monitorChan:
			data,err := monitor.Encode()
			if nil != err {
				log.Error(err)
				return 
			}

			_, err = SendPostRequestHttp(config.AppConf.Monitor,data)
			if nil != err {
				return 
			}
		}
	}
}

func Cross() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin") //请求头部
		if origin != "" {
			fmt.Println("origin-->", origin)
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Origin", origin)                                     // 这是允许访问所有域
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE") //服务器支持的所有跨域请求的方法,为了避免浏览次请求的多次'预检'请求
			//  header的类型
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session,X_Requested_With,Accept, Origin, Host, Connection, Accept-Encoding, Accept-Language,DNT, X-CustomHeader, Keep-Alive, User-Agent, X-Requested-With, If-Modified-Since, Cache-Control, Content-Type, Pragma")
			//              允许跨域设置                                                                                                      可以返回其他子段
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type,Expires,Last-Modified,Pragma,FooBar") // 跨域关键设置 让浏览器可以解析
			c.Header("Access-Control-Max-Age", "172800")                                                                                                                                                           // 缓存请求信息 单位为秒
			c.Header("Access-Control-Allow-Credentials", "true")                                                                                                                                                  //  跨域请求是否需要带cookie信息 默认设置为true
		}
		// 处理请求
		c.Next() //  处理请求
	}
}

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
	r.Use(Cross())
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