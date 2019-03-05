package customserver

import (
	"net"
	"fmt"
	"time"
	"strings"
	"tsstream/config"
	"github.com/gansidui/gotcp"
)

type StreamPacket struct {
	buf []byte
}

func (this* StreamPacket)Serialize() []byte {
	return this.buf
}

func NewStreamPacket(buf []byte) *StreamPacket {
	p := &StreamPacket{}
	if 0 == len(buf) {
		return p
	}

	p.buf = make([]byte,len(buf))
	copy(p.buf,buf)

	return p
}

type StreamCallback struct {
	Conns	map[gotcp.Conn]*gotcp.Conn
}

func GetStreamCallback() *StreamCallback {
	callback := &StreamCallback{Conns: make(map[gotcp.Conn]*gotcp.Conn)}

	return callback
}

func (this* StreamCallback)ValidateAccess(addr string) (bool) {
	if  !config.AuthorityConf.Use {
		return true
	}

	ips := strings.Split(config.AuthorityConf.Ips,",")
	if 0 == len(ips) {
		return true
	}

	for _,ip := range(ips) {
		if addr == ip {
			return true
		}
	}
	
	return false
}

/**
* 使用时用于过滤数据源的连接，不讲数据源的连接放入连接池中
*/
func (this* StreamCallback)ValidateFilter(addr string) (bool) {
	ips := strings.Split(config.CustomConf.FilterIps,",")
	if 0 == len(ips) {
		return true
	}

	for _,ip := range(ips) {
		if addr == ip {
			return true
		}
	}
	
	return false
}

func (this* StreamCallback)OnConnect(conn *gotcp.Conn) bool {
	addr := conn.GetRawConn().RemoteAddr()
	if !this.ValidateAccess(addr.String()) {
		fmt.Println("Access denied! IP: " ,addr.String())
		conn.Close()
		
		return false
	}

	conn.PutExtraData(addr)
	
	// 将连接加入连接池中
	if !this.ValidateFilter(addr.String()) {
		this.Conns[*conn] = conn
	}

	fmt.Println("OnConnect from :",addr)

	return true
}

func (this* StreamCallback)OnMessage(conn *gotcp.Conn,packet gotcp.Packet) bool {

	return true
}

func (this* StreamCallback)OnClose(conn *gotcp.Conn) {

	return 
}

/**
* 将数据广播出去
*/
func (this* StreamCallback)BroasdcastData(data []byte) {
	if 0 == len(data) {
		return 
	}

	for _,conn := range(this.Conns) {
		conn.AsyncWritePacket(NewStreamPacket(data),time.Second * 2)
	} 

	fmt.Println("Broadcast to client by TCP,didn't use dispatch :",len(data))
}

type StreamProtocol struct {

}

/**
* 从TCP连接中获取数据，并用数据初始化一个StreamPacket
*/
func (this* StreamProtocol)ReadPacket(conn *net.TCPConn) (gotcp.Packet,error) {
	buf := make([]byte, RecvBufferSize)
	n, err := conn.Read(buf)
	if nil != err {
		return nil, err
	}

	return NewStreamPacket(buf[:n]),nil
}

/**
* tcp服务实例
*/
var custom_tcp_server *gotcp.Server

type TcpServer struct {
	Host 	string 
	Port 	string

	callback *StreamCallback
}

func GetTcpServer(bind,port string) *TcpServer {
	return &TcpServer{Host: bind,Port: port}
}

func (this *TcpServer)BroadcastLoop() {
	if nil == custom_tcp_server {
		return 
	}

	if nil == this.callback {
		return 
	}

	for {
		select {
		case message := <-WebSocketDataChan:
			this.callback.BroasdcastData(message)
		}
	}
}

/**
* 运行流接收服务器
*/
func (this* TcpServer)RunCustomServer() {
	if nil != custom_tcp_server {
		fmt.Println("custom Tcp Server already Run!")

		return 
	}

	addr := fmt.Sprintf("%s:%s", this.Host,this.Port)
	tcpAddr, err := net.ResolveTCPAddr("tcp4", addr)
	if nil != err {
		fmt.Println(err)

		return 
	}

	listener, err := net.ListenTCP("tcp", tcpAddr)
	if nil != err {
		fmt.Println(err)

		return 
	}

	config := &gotcp.Config{
		PacketSendChanLimit:    PacketSendChanLimit,
		PacketReceiveChanLimit: PacketReceiveChanLimit,
	}

	this.callback = GetStreamCallback()
	custom_tcp_server = gotcp.NewServer(config, this.callback, &StreamProtocol{})
	if nil == custom_tcp_server {
		fmt.Println("Create custom Tcp Server Error")

		return 
	}

	defer func() {
		custom_tcp_server.Stop()
		custom_tcp_server = nil
	}()

	fmt.Println("TCP: RunCustomServer is Running ",addr)

	go this.BroadcastLoop()
	custom_tcp_server.Start(listener, time.Second)

	return 
}