package stream

import (
	"net"
	"fmt"
	"time"
	"strings"
	"tsstream/config"
	"tsstream/dispatch"
	"github.com/gansidui/gotcp"
)

/*
* 主要是用于作为一个分布式的TCP流接收服务，并将流通过分发器分发出去
*/

const (
	PacketSendChanLimit 	= 32	// 传送包的chan大小
	PacketReceiveChanLimit 	= 32	// 接收包的chan限制
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
	dispatch dispatch.Dispatch
}

func GetStreamCallback() *StreamCallback {
	callback := &StreamCallback{}
	callback.dispatch.Init()

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

func (this* StreamCallback)OnConnect(conn *gotcp.Conn) bool {
	addr := conn.GetRawConn().RemoteAddr()
	if !this.ValidateAccess(addr.String()) {
		fmt.Println("Access denied! IP: " ,addr.String())

		return false
	}

	conn.PutExtraData(addr)

	fmt.Println("OnConnect from :",addr)

	return true
}

/**
* 接收消息，并将消息分发出去
*/
func (this* StreamCallback)OnMessage(conn *gotcp.Conn,packet gotcp.Packet) bool {
	streamPacket := packet.(*StreamPacket)
	err := this.dispatch.Dispatch(streamPacket.Serialize())
	if nil != err {
		fmt.Println("Dispatch data Error: ",err)
		
		return false
	}

	return true
}

func (this* StreamCallback)OnClose(conn *gotcp.Conn) {

	return 
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
var push_tcp_server *gotcp.Server

type TcpStream struct {
	Host 	string 
	Port 	string
}

func GetTcpStream(bind,port string) *TcpStream {
	return &TcpStream{Host: bind,Port: port}
}

/**
* 运行流接收服务器
*/
func (this* TcpStream)RunServer() {
	if nil != push_tcp_server {
		fmt.Println("Tcp Server already Run!")

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

	callback := GetStreamCallback()
	push_tcp_server = gotcp.NewServer(config, callback, &StreamProtocol{})
	if nil == push_tcp_server {
		fmt.Println("Create Tcp Server Error")

		return 
	}

	defer func() {
		callback.dispatch.Free()
		push_tcp_server.Stop()
		push_tcp_server = nil
	}()

	fmt.Println("TCP: RunStreamServer is Running ",addr)

	push_tcp_server.Start(listener, time.Second)

	return 
}