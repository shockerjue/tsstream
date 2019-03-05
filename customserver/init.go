package customserver

/**
* 数据分发器
*/
var WebSocketDataChan chan []byte

const (
	ClientChanSize = 32
	WebSocketDataSize = 32

	RecvBufferSize = 64 * 1024
)

func init() {
	WebSocketDataChan = make(chan []byte,WebSocketDataSize)
}