package gnet

// 连接的编解码接口
type Codec interface {
	// 包头长度
	// 应用层可以自己扩展包头长度
	PacketHeaderSize() uint32

	// 编码接口
	Encode(connection Connection, packet Packet) []byte

	// 创建消息头
	// packet可能为nil
	// packetData是packet encode后的数据,可能为nil
	CreatePacketHeader(connection Connection, packet Packet, packetData []byte) PacketHeader

	// 解码接口
	Decode(connection Connection, data []byte) (newPacket Packet, err error)
}
