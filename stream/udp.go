package stream

import (
	"net"
	"fmt"
	"runtime"
	"errors"
	"strings"
	"tsstream/config"
	"tsstream/dispatch"
)

/*
* 主要是用于作为一个分布式的UDP流接收服务，并将流通过分发器分发出去
*/

type UdpStream struct {
	Host 	string 
	Port 	string 
	dispatch dispatch.Dispatch
}

func GetUdpStream(host,port string) *UdpStream {
	udp := &UdpStream{Host: host, Port: port}
	udp.dispatch.Init()

	return udp
}

func (this* UdpStream)ValidateAccess(addr string) (bool) {
	if !config.AuthorityConf.Use {
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

func (this* UdpStream)HandleRequested(conn *net.UDPConn) (err error){
	buffer := make([]byte, RecvBufferSize)
	n, addr, err := conn.ReadFromUDP(buffer)

	if !this.ValidateAccess(addr.IP.String()) {
		e_info := "Access denied! IP: " + addr.IP.String()
		err = errors.New(e_info)

		return 
	}

    if err != nil {
		return 
    } else {
		this.dispatch.Dispatch(buffer[:n])
	}
	
	return nil
}

func (this *UdpStream)FreeDispatch() {
	this.dispatch.Free()
}

func (this *UdpStream)RunServer() {
	addr := fmt.Sprintf("%s:%s", this.Host,this.Port)
	udpAddr, err := net.ResolveUDPAddr("udp4", addr)
    if err != nil {
		fmt.Println(err)

		return 
    }

	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		fmt.Println(err)
		
		return 
	}

	defer func() {
		conn.Close()
		this.FreeDispatch()
	}()
	
	fmt.Println("UDP: RunStreamServer is Running ",addr)

	for {
		err = this.HandleRequested(conn)
		if nil != err {
			fmt.Println("Dispatch data Error :", err)
			
			continue
		}
		runtime.Gosched()
	}
}