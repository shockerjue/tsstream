package controller

import (
	"time"
	"tsstream/config"
	"tsstream/stream"
	"tsstream/monitor"
	"tsstream/customserver"
	log "github.com/sirupsen/logrus"
)

func recvStream() {
	err := stream.RunGetStream()
	if nil != err {
		log.Error(err)
	}
}

func webCustom() {
	/**
	* 启动websocket服务，等待客户端的请求
	*/
	if "WEBSOCKET" == config.CustomConf.Protocol {
		err := customserver.RunWebSocketServer()
		if nil != err {
			log.Error(err)
		}

		return 
	}
	
	/**
	* 启动TCP服务，等待客户端的连接
	*/
	if "TCP" == config.CustomConf.Protocol {
		customserver.RunTcpServer()

		return
	}

	log.Println("start CUSTOM config error!")
}

/**
* 这种运行方式必须是起始节点，也就是接收ffmpeg推送来的流的节点
* 接收以后可以将流数据通过TCP/UDP/CUSTOM分发出去，
* 如果使用CUSTOM分发，那就说明这个节点是起始节点，也是末端节点。
*/
func RunNormal() {
	if "CUSTOM" == config.DispatchConf.Protocol {
		go webCustom()
	}

	time.Sleep(2 * time.Second)

	go recvStream()
}

/**
* 该允许方式是接收来自TCP、UDP的数据，并将数据通过TCP/UDP/CUSTOM分发出去
* 这里可以作为中转节点，或者末端节点，这种方式允许的节点可以串联起来形成一条链式结构
* 如果形成链试结构，那么分发必须通过TCP/UDP来分发。使用CUSTOM分发，那就说运行
* 在末端
*/
func RunExtra() {
	/**
	* 启动端点服务，即启动服务来服务客户请求
	*/
	if "CUSTOM" == config.DispatchConf.Protocol {
		go webCustom()
	}

	time.Sleep(3 * time.Second)

	go stream.RunPushServer()
}

func RunMonitor() {
	monitor := monitor.GetMonitorServer()
	go monitor.RunServer()
}