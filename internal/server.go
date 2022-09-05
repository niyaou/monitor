package internal

import (
	"context"
	"sync"
	"time"

	. "desay.com/radar-monitor/gnet"
	"desay.com/radar-monitor/logger"
)

var (
	// singleton
	_server Server
)

// 服务器接口
type Server interface {
	GetContext() context.Context

	GetWaitGroup() *sync.WaitGroup

	SetBroker(*BrokerClient)

	GetBroker() *BrokerClient

	// 初始化
	Init(ctx context.Context, configFile string) bool

	// 运行、启动链接
	Run(ctx context.Context)

	// 关闭链接
	Stop(ctx context.Context)

	// 定时更新
	OnUpdate(ctx context.Context, updateCount int64)

	// 退出
	Exit()
}

func SetServer(server Server) {
	_server = server
}

func GetServer() Server {
	return _server
}

type BaseServerConfig struct {
	// 服务器id
	ServerId int32

	// 客户端监听地址
	ClientListenAddr string
	// 客户端连接配置
	ClientConnConfig ConnectionConfig
	// 服务器监听地址
	ServerListenAddr string
	// 服务器连接配置
	ServerConnConfig ConnectionConfig
}

// 服务器基础流程
type BaseServer struct {
	// 定时更新间隔
	updateInterval time.Duration
	// 更新次数
	updateCount int64

	configFile string

	ctx context.Context

	wg sync.WaitGroup

	broker *BrokerClient

	onServerInitFuncs []func()
}

func (_self *BaseServer) GetConfigFile() string {
	return _self.configFile
}

func (_self *BaseServer) SetBroker(broker *BrokerClient) {
	_self.broker = broker
}

func (_self *BaseServer) GetBroker() *BrokerClient {
	return _self.broker
}

func (_self *BaseServer) GetContext() context.Context {
	return _self.ctx
}

func (_self *BaseServer) GetWaitGroup() *sync.WaitGroup {
	return &_self.wg
}

// // 添加初始化回调函数
func (_self *BaseServer) AddInitHook(initFunc ...func()) {
	_self.onServerInitFuncs = append(_self.onServerInitFuncs, initFunc...)
}

// 加载配置文件
func (_self *BaseServer) Init(ctx context.Context, configFile string) bool {
	_self.configFile = configFile
	_self.updateInterval = time.Second
	_self.ctx = ctx
	for _, initFunc := range _self.onServerInitFuncs {
		initFunc()
	}
	return true
}

// 运行
func (_self *BaseServer) Run(ctx context.Context) {
	// go func(ctx context.Context) {
	// 	_self.updateLoop(ctx)
	// }(ctx)
}

// 关闭连接
func (_self *BaseServer) Stop(ctx context.Context) {
	// go func(ctx context.Context) {
	// 	_self.updateLoop(ctx)
	// }(ctx)
}

func (_self *BaseServer) OnUpdate(ctx context.Context, updateCount int64) {
}

func (_self *BaseServer) Exit() {
	logger.Info("BaseServer.Exit")
	// 服务器管理的协程关闭
	logger.Info("wait server goroutine close")
	_self.wg.Wait()
	logger.Info("all server goroutine closed")
	// 网络关闭
	logger.Info("wait net goroutine close")
	GetNetMgr().Shutdown(true)
	logger.Info("all net goroutine closed")
	_self.broker.Close()

}

// 定时更新接口
func (_self *BaseServer) updateLoop(ctx context.Context) {
	// 暂定更新间隔1秒
	updateTicker := time.NewTicker(_self.updateInterval * 20)
	defer func() {
		updateTicker.Stop()
		logger.Info("updateLoop end")
	}()
	for {
		select {
		// 系统关闭通知
		case <-ctx.Done():
			logger.Info("exitNotify")
			return
		case <-updateTicker.C:
			_self.OnUpdate(ctx, _self.updateCount)
			_self.updateCount++
		}
	}
}
