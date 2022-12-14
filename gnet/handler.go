package gnet

// 连接回调
type ConnectionHandler interface {
	// 连接成功或失败
	OnConnected(connection Connection, success bool)

	// 断开连接
	OnDisconnected(connection Connection)

	// 收到一个完整数据包
	// 在收包协程中调用
	OnRecvPacket(connection Connection, packet Packet)

	// 创建一个心跳包(只对connector有效)
	// 在connector的发包协程中调用
	CreateHeartBeatPacket(connection Connection) Packet
}

// 监听回调
type ListenerHandler interface {
	// accept a new connection
	OnConnectionConnected(listener Listener, acceptedConnection Connection)

	// a connection disconnect
	OnConnectionDisconnect(listener Listener, connection Connection)
}
