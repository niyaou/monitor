package gnet

import (
	"context"
	"sync"
)

// 网络管理类,提供对外接口
type NetMgr struct {

	// 监听对象管理
	listenerMap     map[uint32]Listener
	listenerMapLock sync.RWMutex

	// 连接对象管理
	connectorMap     map[uint32]Connection
	connectorMapLock sync.RWMutex

	// 初始化一次
	initOnce sync.Once
	// 管理协程的关闭
	wg sync.WaitGroup
}

var (
	// singleton
	netMgr = &NetMgr{}
)

// 单例模式,在调用的时候才会执行初始化一次
func GetNetMgr() *NetMgr {
	netMgr.initOnce.Do(func() {
		netMgr.init()
	})
	return netMgr
}

// 初始化
func (_self *NetMgr) init() {
	_self.listenerMap = make(map[uint32]Listener)
	_self.connectorMap = make(map[uint32]Connection)
	_self.wg = sync.WaitGroup{}
}

// 新监听对象
func (_self *NetMgr) NewListener(ctx context.Context, address string, acceptConnectionConfig ConnectionConfig, acceptConnectionCodec Codec,
	acceptConnectionHandler ConnectionHandler, listenerHandler ListenerHandler) Listener {
	newListener := NewTcpListener(acceptConnectionConfig, acceptConnectionCodec, acceptConnectionHandler, listenerHandler)
	newListener.netMgrWg = &_self.wg
	if !newListener.Start(ctx, address) {
		logger.Debug("NewListener Start Failed")
		return nil
	}
	_self.listenerMapLock.Lock()
	_self.listenerMap[newListener.GetListenerId()] = newListener
	_self.listenerMapLock.Unlock()

	newListener.onClose = func(listener Listener) {
		_self.listenerMapLock.Lock()
		delete(_self.listenerMap, listener.GetListenerId())
		_self.listenerMapLock.Unlock()
	}
	return newListener
}

func (_self *NetMgr) NewListenerCustom(ctx context.Context, address string, acceptConnectionConfig ConnectionConfig, acceptConnectionCodec Codec,
	acceptConnectionHandler ConnectionHandler, listenerHandler ListenerHandler, acceptConnectionCreator AcceptConnectionCreator) Listener {
	newListener := NewTcpListener(acceptConnectionConfig, acceptConnectionCodec, acceptConnectionHandler, listenerHandler)
	newListener.acceptConnectionCreator = acceptConnectionCreator
	newListener.netMgrWg = &_self.wg
	if !newListener.Start(ctx, address) {
		logger.Debug("NewListener Start Failed")
		return nil
	}
	_self.listenerMapLock.Lock()
	_self.listenerMap[newListener.GetListenerId()] = newListener
	_self.listenerMapLock.Unlock()

	newListener.onClose = func(listener Listener) {
		_self.listenerMapLock.Lock()
		delete(_self.listenerMap, listener.GetListenerId())
		_self.listenerMapLock.Unlock()
	}
	return newListener
}

// 新连接对象
func (_self *NetMgr) NewConnector(ctx context.Context, address string, connectionConfig *ConnectionConfig,
	codec Codec, handler ConnectionHandler, tag interface{}) Connection {
	return _self.NewConnectorCustom(ctx, address, connectionConfig, codec, handler, tag, func(_config *ConnectionConfig, _address string, _codec Codec, _handler ConnectionHandler, isFixedHead bool) Connection {
		return NewTcpConnector(_config, _address, _codec, _handler, false)
	}, false)
}

func (_self *NetMgr) NewConnectorCustom(ctx context.Context, address string, connectionConfig *ConnectionConfig,
	codec Codec, handler ConnectionHandler, tag interface{}, connectionCreator ConnectionCreator, isFixedHead bool) Connection {
	newConnector := connectionCreator(connectionConfig, address, codec, handler, isFixedHead)
	newConnector.SetAddr(address)
	newConnector.SetTag(tag)
	// if !newConnector.Connect(address) {
	// 	newConnector.Close()
	// 	return nil
	// }
	_self.connectorMapLock.Lock()
	_self.connectorMap[newConnector.GetConnectionId()] = newConnector
	_self.connectorMapLock.Unlock()
	// newConnector.Start(ctx, &_self.wg, func(connection Connection) {
	// 	_self.connectorMapLock.Lock()
	// 	delete(_self.connectorMap, connection.GetConnectionId())
	// 	_self.connectorMapLock.Unlock()
	// })
	return newConnector
}

func (_self *NetMgr) StartConnector(ctx context.Context, address string, newConnector Connection) bool {
	if !newConnector.Connect(address) {
		newConnector.Close()
		return false
	}
	newConnector.Start(ctx, &_self.wg, func(connection Connection) {
		_self.connectorMapLock.Lock()
		delete(_self.connectorMap, connection.GetConnectionId())
		_self.connectorMapLock.Unlock()
	})
	return true
}

func (_self *NetMgr) StopConnector(ctx context.Context, conn Connection) bool {
	return false
}

// 关闭
// waitForAllNetGoroutine:是否阻塞等待所有网络协程结束
func (_self *NetMgr) Shutdown(waitForAllNetGoroutine bool) {
	if waitForAllNetGoroutine {
		// 等待所有网络协程结束
		_self.wg.Wait()
		logger.Debug("all net goroutine closed")
	}
}
