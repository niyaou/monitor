package gnet

// ProtoPacket默认ConnectionHandler
type DefaultConnectionHandler struct {
	// 注册消息的处理函数map
	PacketHandlers map[PacketCommand]ProtoPacketHandler
	// 未注册消息的处理函数
	UnRegisterHandler ProtoPacketHandler
	// 连接回调
	onConnectedFunc    func(connection Connection, success bool)
	onDisconnectedFunc func(connection Connection)
	// handler一般总是和codec配合使用
	protoCodec Codec
	// 心跳包消息号(只对connector有效)
	heartBeatCommand PacketCommand
	// 心跳包构造函数(只对connector有效)
	heartBeatCreator ProtoMessageCreator
}
type PacketHandlerRegister interface {
	Register(packetCommand PacketCommand, handler ProtoPacketHandler, creator ProtoMessageCreator)
}

// ProtoPacket消息回调
type ProtoPacketHandler func(connection Connection, packet *ProtoPacket)

func (_self *DefaultConnectionHandler) OnConnected(connection Connection, success bool) {
	if _self.onConnectedFunc != nil {
		_self.onConnectedFunc(connection, success)
	}
}

func (_self *DefaultConnectionHandler) OnDisconnected(connection Connection) {
	if _self.onDisconnectedFunc != nil {
		_self.onDisconnectedFunc(connection)
	}
}

func (_self *DefaultConnectionHandler) OnRecvPacket(connection Connection, packet Packet) {
	// logger.Info(" OnRecvPacket:%v", packet)
	defer func() {
		if err := recover(); err != nil {
			logger.Error("fatal %v", err.(error))
			LogStack()
		}
	}()
	if protoPacket, ok := packet.(*ProtoPacket); ok {
		if packetHandler, ok2 := _self.PacketHandlers[protoPacket.command]; ok2 {
			if packetHandler != nil {
				packetHandler(connection, protoPacket)
				return
			}
		}
		if _self.UnRegisterHandler != nil {
			_self.UnRegisterHandler(connection, protoPacket)
		}
	}
}

func (_self *DefaultConnectionHandler) CreateHeartBeatPacket(connection Connection) Packet {
	if _self.heartBeatCreator != nil {
		return NewProtoPacket(_self.heartBeatCommand, _self.heartBeatCreator())
	}
	return nil
}

func NewDefaultConnectionHandler(protoCodec Codec) *DefaultConnectionHandler {
	return &DefaultConnectionHandler{
		PacketHandlers: make(map[PacketCommand]ProtoPacketHandler),
		protoCodec:     protoCodec,
	}
}

func (_self *DefaultConnectionHandler) GetCodec() Codec {
	return _self.protoCodec
}

// 注册消息号和消息回调,消息构造的映射
// handler在TcpConnection的read协程中被调用
func (_self *DefaultConnectionHandler) Register(packetCommand PacketCommand, handler ProtoPacketHandler, creator ProtoMessageCreator) {
	_self.PacketHandlers[packetCommand] = handler
	if _self.protoCodec != nil && creator != nil {
		if protoRegister, ok := _self.protoCodec.(ProtoRegister); ok {
			protoRegister.Register(packetCommand, creator)
		}
	}
}

func (_self *DefaultConnectionHandler) GetPacketHandler(packetCommand PacketCommand) ProtoPacketHandler {
	return _self.PacketHandlers[packetCommand]
}

// 注册心跳包(只对connector有效)
func (_self *DefaultConnectionHandler) RegisterHeartBeat(packetCommand PacketCommand, creator ProtoMessageCreator) {
	_self.heartBeatCommand = packetCommand
	_self.heartBeatCreator = creator
}

// 未注册消息的处理函数
// unRegisterHandler在TcpConnection的read协程中被调用
func (_self *DefaultConnectionHandler) SetUnRegisterHandler(unRegisterHandler ProtoPacketHandler) {
	_self.UnRegisterHandler = unRegisterHandler
}

// 设置连接回调
func (_self *DefaultConnectionHandler) SetOnConnectedFunc(onConnectedFunc func(connection Connection, success bool)) {
	_self.onConnectedFunc = onConnectedFunc
}

// 设置连接断开回调
func (_self *DefaultConnectionHandler) SetOnDisconnectedFunc(onDisconnectedFunc func(connection Connection)) {
	_self.onDisconnectedFunc = onDisconnectedFunc
}
