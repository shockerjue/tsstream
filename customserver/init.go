package customserver

const (
	ClientChanSize = 32
	WebSocketDataSize = 32
	RecvBufferSize = 64 * 1024
	PacketSendChanLimit 	= 128	// 传送包的chan大小
	PacketReceiveChanLimit 	= 128	// 接收包的chan限制
)

/**
* 数据分发器
*/
var WebSocketDataChan chan []byte
func init() {
	WebSocketDataChan = make(chan []byte,WebSocketDataSize)
}