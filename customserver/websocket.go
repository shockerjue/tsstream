package customserver 

import (
	"runtime"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

type Client struct {
	hub 		*WebSocketServer
	conn 		*websocket.Conn
	send 		chan []byte
}

/**
* 在协程中发送消息
*/
func (this *Client) writePump() {
	defer func() {
		this.hub.unregister <- this
		this.conn.Close()
		log.Error("Client connect close!")
	}()

	for {
		select {
			case message,ok := <- this.send:
				if !ok {
					log.Println("WebSocket writePump over? ok: ",ok)

					return 
				}

				/**
				* 需要设置为发送二进制数据，否则浏览器会报出UTF-8解码错误
				*/
				w,err := this.conn.NextWriter(websocket.BinaryMessage)
				if err != nil {
					log.Println("WebSocket writePump over? NextWriter: ",err)

					return 
				}

				w.Write(message)
				if err := w.Close(); err != nil {
					log.Println("WebSocket writePump over? Close: ",err)

					return 
				}

				n := len(this.send)
				for i := 0; i < n - 1; i++ {
					data := <-this.send
					w,err := this.conn.NextWriter(websocket.BinaryMessage)
					if err != nil {
						log.Println("WebSocket writePump over? NextWriter: ",err)

						continue 
					}

					runtime.Gosched()
					w.Write(data)

					if err := w.Close(); err != nil {
						log.Println("WebSocket writePump over? Close: ",err)
	
						continue 
					}
				}
		}
	}

}

func WebSocketConn(hub *WebSocketServer,w http.ResponseWriter,r *http.Request,c *gin.Context) {
	var upgrade = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},

		ReadBufferSize: 128,
		WriteBufferSize: 64 * 1024 * 1,
	}

	/**
	* 设置响应头部，要不会报错，也就是不设置就连接不上
	*/
	var headers http.Header = make(http.Header)
	headers.Add("Sec-WebSocket-Protocol","null")

	conn, err := upgrade.Upgrade(w,r,headers) 
	if nil != err {
		log.Println(err)

		return 
	}

	client := &Client{hub: hub,conn: conn,send: make(chan []byte,ClientChanSize)}
	client.hub.register <- client

	go client.writePump()
}


type WebSocketServer struct {
	clients map[*Client]bool
	register chan *Client
	unregister chan *Client
}

func NewWebSocketServer() *WebSocketServer {
	return &WebSocketServer{
		register: make(chan *Client),
		unregister: make(chan *Client),
		clients: make(map[*Client]bool),
	}
}

func (h *WebSocketServer) RunServer() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true

		case client := <-h.unregister:
			if _,ok := h.clients[client]; ok {
				delete(h.clients,client)
				close(client.send)
				log.Error("Client chan close!")
			}

		case message := <- WebSocketDataChan:
			for client := range h.clients {
				if (ClientChanSize - 10) < len(client.send) {
					continue
				}
				 
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients,client)
				}
			}

			n := len(WebSocketDataChan)
			for i := 0; i < n - 1; i++ {
				data := <-WebSocketDataChan
				for client := range h.clients {
					runtime.Gosched()
					if (ClientChanSize - 10) < len(client.send) {
						continue
					}

					select {
					case client.send <- data:
					default:
						close(client.send)
						delete(h.clients,client)
					}
				}
			}
		}
	}
}