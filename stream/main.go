package stream

import (
	"fmt"
	"tsstream/config"
	log "github.com/sirupsen/logrus"
)

/**
* 启动UDP数据流接收服务
* 主要是接受来自ffmpeg流的数据
*/
func udpStream() (err error) {
	host := config.BackStreamConf.Bind
	port := config.BackStreamConf.Port
	
	udpStream := GetUdpStream(host,port)
	udpStream.RunServer()

	return 
}

/**
* 启动TCP数据流接收服务
* 接受来自ffmpeg流的数据
*/
func tcpStream() (err error) {
	host := config.BackStreamConf.Bind
	port := config.BackStreamConf.Port

	tcpStream := GetTcpStream(host,port)
	tcpStream.RunServer()

	return 
}

/**
* 启动http数据流接收服务
*/
func httpStream() (err error) {
	host := config.BackStreamConf.Bind
	port := config.BackStreamConf.Port

	httpStream := GetHttpStream(host,port)

	httpStream.RunServer()

	return 
}

/**
* 根据接受流的协议方式运行服务
* 这里主要是接受来自ffmpeg推上来的数据，并根据分发方式将流分发出去
* 分发的目标包含分发到其他节点、进程、websocket服务
*/
func RunGetStream() (err error){
	protocol := config.BackStreamConf.Protocol
	fmt.Println("Start Get Stream Server, protocol:",protocol)

	if "" == protocol {
		protocol = "HTTP"
	}

	if "HTTP" == protocol {
		err = httpStream()

		return 
	} 
	
	if "TCP" == protocol {
		err = tcpStream()

		return 
	}
	
	if "UDP" == protocol {
		err = udpStream()

		return 
	}

	fmt.Println("Didn't support stream protocol")

	return 
}	

/**
* 根据接收推流的协议允许接受流的服务
* 主要是用于分布式，或者不通进程接受流，并将其分发到客户端中
* 服务主要包含UDP服务或者TCP服务
*/
func RunPushServer() (err error) {
	bind := config.PushStreamConf.Bind
	port := config.PushStreamConf.Port
	protocol := config.PushStreamConf.Protocol

	if "TCP" == protocol {
		tcpServer := GetTcpStream(bind,port)
		go tcpServer.RunServer()
	}

	if "UDP" == protocol {
		udpServer := GetUdpStream(bind,port)
		go udpServer.RunServer()
	}

	log.Println("Start Push Server ",protocol,"  ",bind,":",port)

	return 
}