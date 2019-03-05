package config

import (
	"log"
	"github.com/go-ini/ini"
)

var (
	AppConf = &AppConfig{}
	BackStreamConf = &BackStreamConfig{}
	PushStreamConf = &PushStreamConfig{}
	DispatchConf = &DispatchConfig{}
	AuthorityConf = &AuthorityConfig{}
	CustomConf = &CustomConfig{}
)

/**
* 运行模式
* normal 	
* 	主要是作为正常节点启动，这种方式启动主要是获取来自ffmpeg推送的流
* 	根据分发器将流分发出去，分发到其他节点(UDP/TCP)或者是客户服务(WEBSOCKET/TCP)
*
* extra
* 	主要是作为扩展服务启动，这种方式主要是接收来自normal启动的服务推送过来的数据
* 	根据分发器将数据分发出去，分发到其他节点(TCP/UDP)或者是客户服务(WEBSOCKET/TCP)
*/
type AppConfig struct {
	RunMode	string
}

/**
* 接收来自ffmpeg推流过来的服务配置信息
* ffmpeg推流的协议与这个配置有关
* 这里的配置信息主要是ffmpeg推流直接交互
* protocol (TCP/UDP/HTTP)
*/
type BackStreamConfig struct {
	Protocol 	string		//	接收流的协议，并且根据这个创建对应的服务(TCP/UDP/HTTP)
	Bind 		string		//	服务绑定的信息
	Port 		string		//	
	Secret		string		//	秘钥（主要是使用HTTP协议时使用）
}

/**
* 接收通过TCP/UDP推流过来的服务配置信息
* 这里主要是接收通过TCP/UDP推送过来的数据，并将其通过dispatch分发出去
* 这里根据下面的配置创建接收数据的服务
* protocol (TCP/UDP)
*/
type PushStreamConfig struct {
	Protocol		string 	//	接收流的协议
	Bind 			string	//	绑定的地址
	Port 			string	//	绑定的端口
}

/**
* 分发器信息配置，接收到数据以后使用分发器将数据分发到其他节点、或者客户
* protocol (TCP/UDP/CUSTOM)
*/
type DispatchConfig struct {
	Hosts		string 		//	分发的主机列表
	Ports		string 		//	分发的端口列表
	Protocol	string		//	分发的协议 （TCP/UDP/CUSTOM）
}

/**
* 授权配置,一般是授权推流过来的连接
* 要使用授权设置，需要将Use设置为true，
* 如果使用了，则仅仅接受IPS的数据
*/
type AuthorityConfig struct {
	Use 	bool
	Ips 	string
}

/**
* 客户服务配置信息，只有在分发服务中使用CUSTOM协议类型这里的配置才有用
* protocol (WEBSOCKET/TCP)
* WEBSOCKET
* 		客户连接使用的是websocket连接服务
*
* TCP 
* 		客户使用的是TCP连接服务,且使用TCP时，
* 		需要过滤推流过来的连接，过滤的连接在FilterIps中设置
*/
type CustomConfig struct {
	Protocol	string		//	客户端使用的协议
	FilterIps	string		//	当使用TCP时需要过滤的IP地址
	Bind 		string		//	客户服务绑定的地址
	Port 		string		//	客户服务绑定的端口
}

func init() {
	cfg, err := ini.Load("conf/app.conf")
	if err != nil {
		log.Fatal("load app conf err:", err)
	}

	err = cfg.Section("app").MapTo(AppConf)
	if err != nil {
		log.Fatal("init app conf err:", err)
	}

	err = cfg.Section("backstream").MapTo(BackStreamConf)
	if err != nil {
		log.Fatal("init backstream conf err:", err)
	}

	err = cfg.Section("pushstream").MapTo(PushStreamConf)
	if err != nil {
		log.Fatal("init pushstream conf err:", err)
	}

	err = cfg.Section("dispatch").MapTo(DispatchConf)
	if err != nil {
		log.Fatal("init dispatch conf err:", err)
	}

	err = cfg.Section("authority").MapTo(AuthorityConf)
	if err != nil {
		log.Fatal("init authority conf err:", err)
	}
	
	err = cfg.Section("custom").MapTo(CustomConf)
	if err != nil {
		log.Fatal("init custom conf err:", err)
	}

	return
}