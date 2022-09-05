package gnet

import (
	"context"
	"net"
	"sync"
	"syscall"
	"time"
)

// TCP监听
type TcpListener struct {
	baseListener

	netListener             net.Listener
	acceptConnectionConfig  ConnectionConfig
	acceptConnectionCodec   Codec
	acceptConnectionHandler ConnectionHandler

	// 连接表
	connectionMap     map[uint32]*TcpConnection
	connectionMapLock sync.RWMutex

	isRunning bool
	// 防止执行多次关闭操作
	closeOnce sync.Once
	// 关闭回调
	onClose func(listener Listener)

	acceptConnectionCreator AcceptConnectionCreator
	// 外部传进来的WaitGroup
	netMgrWg *sync.WaitGroup
}

func NewTcpListener(acceptConnectionConfig ConnectionConfig, acceptConnectionCodec Codec, acceptConnectionHandler ConnectionHandler, listenerHandler ListenerHandler) *TcpListener {
	return &TcpListener{
		baseListener: baseListener{
			listenerId: newListenerId(),
			handler:    listenerHandler,
		},
		acceptConnectionConfig:  acceptConnectionConfig,
		acceptConnectionCodec:   acceptConnectionCodec,
		acceptConnectionHandler: acceptConnectionHandler,
		connectionMap:           make(map[uint32]*TcpConnection),
	}
}

func (_self *TcpListener) GetConnection(connectionId uint32) Connection {
	_self.connectionMapLock.RLock()
	conn := _self.connectionMap[connectionId]
	_self.connectionMapLock.RUnlock()
	return conn
}

// 广播消息
func (_self *TcpListener) Broadcast(packet Packet) {
	_self.connectionMapLock.RLock()
	for _, conn := range _self.connectionMap {
		if conn.isConnected {
			conn.SendPacket(packet.Clone())
		}
	}
	_self.connectionMapLock.RUnlock()
}

// 开启监听
func (_self *TcpListener) Start(ctx context.Context, listenAddress string) bool {
	var err error
	_self.netListener, err = net.Listen("tcp", listenAddress)
	if err != nil {
		logger.Error("Listen Failed %v: %v", _self.GetListenerId(), err)
		return false
	}

	// 监听协程
	_self.isRunning = true
	_self.netMgrWg.Add(1)
	go func(ctx context.Context) {
		defer _self.netMgrWg.Done()
		_self.acceptLoop(ctx)
	}(ctx)

	// 关闭响应协程
	_self.netMgrWg.Add(1)
	go func() {
		defer _self.netMgrWg.Done()
		select {
		// 关闭通知
		case <-ctx.Done():
			logger.Debug("recv closeNotify %v", _self.GetListenerId())
			_self.Close()
		}
	}()

	return true
}

// 关闭监听,并关闭管理的连接
func (_self *TcpListener) Close() {
	_self.closeOnce.Do(func() {
		_self.isRunning = false
		if _self.netListener != nil {
			_self.netListener.Close()
		}
		connMap := make(map[uint32]*TcpConnection)
		_self.connectionMapLock.RLock()
		for k, v := range _self.connectionMap {
			connMap[k] = v
		}
		_self.connectionMapLock.RUnlock()
		// 关闭管理的连接
		for _, conn := range connMap {
			conn.Close()
		}
		if _self.onClose != nil {
			_self.onClose(_self)
		}
	})
}

// accept协程
func (_self *TcpListener) acceptLoop(ctx context.Context) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("acceptLoop fatal %v: %v", _self.GetListenerId(), err.(error))
			LogStack()
		}
	}()

	for _self.isRunning {
		// 阻塞accept,当netListener关闭时,会返回err
		newConn, err := _self.netListener.Accept()
		if err != nil {
			logger.Error("%v accept err:%v", _self.GetListenerId(), err)
			// 有可能是因为open file数量限制 而导致的accept失败
			if err == syscall.EMFILE {
				logger.Error("accept failed id:%v syscall.EMFILE", _self.GetListenerId())
				// 这个错误只是导致新连接暂时无法连接,不应该退出监听,当有连接释放后,新连接又可以连接上
				time.Sleep(time.Millisecond)
				continue
			}
			break
		}
		_self.netMgrWg.Add(1)
		go func() {
			defer func() {
				_self.netMgrWg.Done()
				if err := recover(); err != nil {
					logger.Error("acceptLoop fatal %v: %v", _self.GetListenerId(), err.(error))
					LogStack()
				}
			}()
			newTcpConn := NewTcpConnectionAccept(newConn, &_self.acceptConnectionConfig, _self.acceptConnectionCodec, _self.acceptConnectionHandler)
			newTcpConn.isConnected = true
			if newTcpConn.handler != nil {
				newTcpConn.handler.OnConnected(newTcpConn, true)
			}
			_self.connectionMapLock.Lock()
			_self.connectionMap[newTcpConn.GetConnectionId()] = newTcpConn
			_self.connectionMapLock.Unlock()
			// newTcpConn.netMgrWg = _self.netMgrWg
			newTcpConn.Start(ctx, _self.netMgrWg, func(connection Connection) {
				if _self.handler != nil {
					_self.handler.OnConnectionDisconnect(_self, connection)
				}
			})

			if _self.handler != nil {
				_self.handler.OnConnectionConnected(_self, newTcpConn)
				newTcpConn.onClose = func(connection Connection) {
					_self.handler.OnConnectionDisconnect(_self, connection)
				}
			}
		}()
	}
}

// Addr returns the listener's network address.
func (_self *TcpListener) Addr() net.Addr {
	if _self.netListener == nil {
		return nil
	}
	return _self.netListener.Addr()
}
