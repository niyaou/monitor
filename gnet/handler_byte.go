package gnet

// ProtoPacket默认ConnectionHandler
type ByteConnectionHandler struct {
	// 注册消息的处理函数map
	// PacketHandlers map[PacketCommand]PacketHandler

	// 连接回调
	onConnectedFunc    func(connection Connection, success bool)
	onDisconnectedFunc func(connection Connection)
	// handler一般总是和codec配合使用
	dataCodec         Codec
	dataPacketHandler DataPacketHandler
}

type DataPacketHandler func(connection Connection, packet *DataPacket)

func (_self *ByteConnectionHandler) OnConnected(connection Connection, success bool) {
	if _self.onConnectedFunc != nil {
		_self.onConnectedFunc(connection, success)
	}
}

func (_self *ByteConnectionHandler) OnDisconnected(connection Connection) {
	if _self.onDisconnectedFunc != nil {
		_self.onDisconnectedFunc(connection)
	}
}

func (_self *ByteConnectionHandler) OnRecvPacket(connection Connection, packet Packet) {
	// logger.Info(" OnRecvPacket:%v", packet)
	defer func() {
		if err := recover(); err != nil {
			logger.Error("fatal %v", err.(error))
			LogStack()
		}
	}()
	if dataPacket, ok := packet.(*DataPacket); ok {

		if _self.dataPacketHandler != nil {
			_self.dataPacketHandler(connection, dataPacket)
			return
		}

		// if _self.UnRegisterHandler != nil {
		// 	_self.UnRegisterHandler(connection, dataPacket)
		// }
	}
}

func NewByteConnectionHandler(dataCodec Codec) *ByteConnectionHandler {
	return &ByteConnectionHandler{
		dataCodec: dataCodec,
	}
}

func (_self *ByteConnectionHandler) GetCodec() Codec {
	return _self.dataCodec
}

// 注册消息号和消息回调,消息构造的映射
// handler在TcpConnection的read协程中被调用
func (_self *ByteConnectionHandler) Register(packetCommand PacketCommand, handler DataPacketHandler) {
	_self.dataPacketHandler = handler

}

func (_self *ByteConnectionHandler) GetPacketHandler(packetCommand PacketCommand) DataPacketHandler {
	return _self.dataPacketHandler
}

// 设置连接回调
func (_self *ByteConnectionHandler) SetOnConnectedFunc(onConnectedFunc func(connection Connection, success bool)) {
	_self.onConnectedFunc = onConnectedFunc
}

// 设置连接断开回调
func (_self *ByteConnectionHandler) SetOnDisconnectedFunc(onDisconnectedFunc func(connection Connection)) {
	_self.onDisconnectedFunc = onDisconnectedFunc
}
func (_self *ByteConnectionHandler) CreateHeartBeatPacket(connection Connection) Packet {
	return nil
}
