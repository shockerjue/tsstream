package dispatch

import (
	"net"
	"fmt"
	"time"
	"errors"
	"strings"
	"tsstream/config"
	"tsstream/customserver"
)

/**
* 分发器，主要是将流通过定义好的书籍分发出去
*/
type Dispatch struct {
	udp 	[]net.Conn	// 分发的UDP连接列表
	tcp 	[]net.Conn	// 分发的TCP连接列表
}

func GetDispatch() *Dispatch {
	return &Dispatch{}
}

/**
* 根据分发协议类型初始化分发连接器
*/
func (this* Dispatch)Init() (err error){
	if "UDP" == config.DispatchConf.Protocol {
		err = this.InitUdp()

		return 
	}

	if "TCP" == config.DispatchConf.Protocol {
		err = this.InitTcp()

		return 
	}

	if "CUSTOM" == config.DispatchConf.Protocol {
		return
	}

	return 
}

/**
*  初始化分发数据流的udp连接列表信息
*/
func (this* Dispatch)InitUdp() (err error){
	this.udp = make([]net.Conn,0)
	hosts := strings.Split(config.DispatchConf.Hosts,",")
	ports := strings.Split(config.DispatchConf.Ports,",")

	if len(ports) != len(hosts) {
		err = errors.New("UDP Hosts != Ports")

		return 
	}

	for k,_ := range(hosts) {
		addr := fmt.Sprintf("%s:%s", hosts[k],ports[k])
		udp,err := net.Dial("udp", addr)
		if nil != err {
			fmt.Println("UDP Connect ",addr," Error ,",err)

			continue 
		}

		this.udp = append(this.udp,udp)

		fmt.Println("UDP Connect ",addr," Success")
	}
	
	return 
}

/**
* 初始化分发流的TCP连接列表信息
*/
func (this* Dispatch)InitTcp() (err error) {
	this.tcp = make([]net.Conn,0)
	hosts := strings.Split(config.DispatchConf.Hosts,",")
	ports := strings.Split(config.DispatchConf.Ports,",")

	if len(ports) != len(hosts) {
		err = errors.New("TCP Hosts != Ports")

		return 
	}

	for k,_ := range(hosts) {
		addr := fmt.Sprintf("%s:%s", hosts[k],ports[k])
		tcp, err := net.DialTimeout("tcp", addr, 2 * time.Second)
		if nil != err {
			fmt.Println("TCP Connect ",addr, " Error,",err)

			continue
		}

		this.tcp = append(this.tcp,tcp)

		fmt.Println("TCP Connect ",addr, " Success")
	}

	return 
}

func (this *Dispatch)Free() {
	if 0 < len(this.udp) {
		for _,conn := range(this.udp) {
			conn.Close()
		}
	}

	if 0 < len(this.tcp) {
		for _,conn := range(this.tcp) {
			conn.Close()
		}
	}

	fmt.Println("Free dispatch")
}

/**
* 将数据分发到UDP服务中，通过初始化的UDP连接列表
*
* @param data 	要分发的数据
*/
func (this* Dispatch)dispatchUdp(data []byte) (err error) {
	if 0 == len(this.udp) {
		err = errors.New("Udp didn't init!")

		return 
	}

	if 0 == len(data) {
		err = errors.New("Udp didn't data!")

		return 
	}

	for _,conn := range(this.udp) {
		conn.Write(data)
	} 

	return
}

/**
* 将数据分发到TCP服务器中，通过初始化的TCP连接列表
*
* @param data 	要分发的数据
*/
func (this* Dispatch)dispatchTcp(data []byte) (err error) {
	if 0 == len(this.tcp) {
		err = errors.New("TCP didn't data!")

		return 
	}

	if 0 == len(data) {
		err = errors.New("TCP didn't data!")

		return 
	}

	for _,conn := range(this.tcp) {
		conn.Write(data)
	}

	return 
}

/**
* 将数据分发到websocket中，它仅仅限于当前进程中
*
* @param data 	要分发的数据
*/
func (this* Dispatch)dispatchWebsocket(data []byte) (err error) {
	if 0 == len(data) {
		return 
	}

	if (customserver.WebSocketDataSize - 10) < len(customserver.WebSocketDataChan) {
		err = errors.New("PushDataChan is filled full!")

		return 
	}
	

	customserver.WebSocketDataChan <- data

	return 
}

/**
* 主要用于将收到根据分发协议分发出去
*
* @param data 
*/
func (this* Dispatch)Dispatch(data []byte) (err error){
	if 0 == len(data) {
		return 
	}

	if "UDP" == config.DispatchConf.Protocol {
		err = this.dispatchUdp(data)

		return 
	}

	if "TCP" == config.DispatchConf.Protocol {
		err = this.dispatchTcp(data)

		return
	}

	/**
	* 支持本地的WEBSOCKET和TCP分发
	*/
	if "CUSTOM" == config.DispatchConf.Protocol {
		err = this.dispatchWebsocket(data)

		return 
	}

	return 
}