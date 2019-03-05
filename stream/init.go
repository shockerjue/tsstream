package stream

const (
	RecvBufferSize = 64 * 1024
	PacketSendChanLimit 	= 32	// 传送包的chan大小
	PacketReceiveChanLimit 	= 32	// 接收包的chan限制
)