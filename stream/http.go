package stream

import (
	"net/http"
	"errors"
	"runtime"
	"tsstream/config"
	"tsstream/dispatch"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

/**
* 接受HTTP流服务处理
*/

type HttpStream struct {
	Host 	string
	Port 	string
	dispatch dispatch.Dispatch
}

func GetHttpStream(host,port string) *HttpStream {
	http := &HttpStream{Host: host, Port: port}
	http.dispatch.Init()

	return http
}

func (this *HttpStream)FreeDispatch() {
	this.dispatch.Free()
}

func (this* HttpStream)HandleRequested(c *gin.Context) (err error) {
	secret := c.Param("secret")
	if secret != config.BackStreamConf.Secret {
		err = errors.New("Secret error!")

		return
	}

	for true {
		buf := make([]byte,RecvBufferSize)
		n, err := c.Request.Body.Read(buf)
		if nil != err || 0 == n {
			break 
		}
		runtime.Gosched()
		err = this.dispatch.Dispatch(buf[:n])
	}

	return 
}

func (this *HttpStream)RunServer() (err error) {
	r := gin.Default()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())


	defer this.FreeDispatch()

	authorized := r.Group("/")
	authorized.Use(func(c *gin.Context) {	
	})
	{
		authorized.POST("/:secret",func(c *gin.Context) {
			err := this.HandleRequested(c)
			if nil != err {
				log.Error("Recv HTTP data error!",err)

				return 
			}

			return 
		})
	}

	s := &http.Server{
		Addr:           ":" + config.BackStreamConf.Port,
		Handler:        r,
		// ReadTimeout:    2 * time.Second, // 移除超时，这样可以连续不断的读取数据,//
		WriteTimeout:   0,
		MaxHeaderBytes: 1 << 20,
	}

	log.Println("Start HTTP Server In ",config.BackStreamConf.Bind," : ",config.BackStreamConf.Port)

	s.ListenAndServe()

	return 
}