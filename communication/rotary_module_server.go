package communication

import (
	"context"
	"encoding/json"
	"os"

	. "desay.com/radar-monitor/gnet"
	. "desay.com/radar-monitor/internal"
	"desay.com/radar-monitor/logger"
)

var (
	_                           Server = (*RotaryModuleServer)(nil)
	_rotaryModuleHandlerChannel        = []string{"ROTARY_MODULE_SEND", "ROTARY_MODULE_ACK"}
	_rotaryModuleMsgMap                = make(map[string][]byte)
	// singleton
	_rotaryModuleServer *RotaryModuleServer
)

type RotaryModuleServer struct {
	BaseServer
	config *RotaryServerConfig
	conn   Connection
	netMgr *NetMgr
}

type RotaryModuleServerConfig struct {
	BaseServerConfig
}

func rotaryModuleSendChannel() {
	s := _rotaryModuleServer
	b := s.BaseServer.GetBroker()
	if b == nil {
		return
	}
	for _, key := range _rotaryModuleHandlerChannel {

		ch, err := b.Subscribe(key)
		if err != nil {
			logger.Error("failed to subscrib: %v", err)
		} else {
			// logger.Info("subscrib channel %v  successed     ", key)
		}

		go func(ch <-chan interface{}, key string) {
			for {
				_msg, ok := _rotaryModuleMsgMap[key]
				if !ok {
					_msg = []byte{}
				}
				_payl := b.GetPayLoad(ch)
				_msg = _payl.([]byte)
				_str_cmd := string(_msg)

				switch key {
				case _rotaryModuleHandlerChannel[0]:
					logger.Info("send-------to channel  to device ---------%v str:%v  ", _msg, _str_cmd)

					if s.conn != nil {

						if s.conn.IsConnected() {

							s.conn.SendPacket(NewDataPacket(nil, _str_cmd))
						} else {

							isConnect := s.netMgr.StartConnector(s.BaseServer.GetContext(), s.conn.RemoteAddr(), s.conn)
							if isConnect {
								s.conn.SendPacket(NewDataPacket(nil, _str_cmd))
							}
						}

					}
				case _rotaryModuleHandlerChannel[1]:
					break
				}

			}
		}(ch, key)
	}
}

func onRtyConnected(connection Connection, success bool) {
	if success {
		_rotaryModuleServer.conn = connection
		// _rotaryModuleServer.conn.SendPacket(NewDataPacket(nil, "RED_angle"))
		logger.Info("------------------------onRtyConnectedt>>>>>>>>>> :%v", "RED_angle")

	}

}

// 初始化
func (_self *RotaryModuleServer) Init(ctx context.Context, configFile string) bool {
	_rotaryModuleServer = _self

	_self.BaseServer.AddInitHook(rotaryModuleSendChannel)

	if !_self.BaseServer.Init(ctx, configFile) {
		return false
	}
	_self.readConfig()

	_self.netMgr = GetNetMgr()

	// clientCodec := NewProtoCodec(nil)
	clientCodec := NewRotaryModuleCodec(nil)
	// clientCodec := NewRotaryModuleCodec(nil)

	clientHandler := NewByteConnectionHandler(clientCodec)

	clientHandler.SetOnConnectedFunc(onRtyConnected)
	_self.registerClientPacket(clientHandler)
	if ctx == nil {
		logger.Info("ctx")
	}
	if _self.config == nil {
		logger.Info("_self.config")
	}
	if clientCodec == nil {
		logger.Info("clientCodec")
	}
	if clientHandler == nil {
		logger.Info("clientHandler")
	}

	_rotaryModuleServer.conn = _self.netMgr.NewConnectorCustom(ctx, _self.config.ClientListenAddr, &_self.config.ClientConnConfig, clientCodec,
		clientHandler, nil, func(config *ConnectionConfig, addr string, codec Codec, handler ConnectionHandler, isFixedHead bool) Connection {
			return NewTcpConnector(config, addr, codec, handler, false)
		}, false)
	if _rotaryModuleServer.conn == nil {
		panic("connect failed")
	}

	return true
}

// 运行
func (_self *RotaryModuleServer) Run(ctx context.Context) {
	_self.BaseServer.Run(ctx)
}

// 退出
func (_self *RotaryModuleServer) Exit() {
	_self.BaseServer.Exit()
}

// 读取配置文件
func (_self *RotaryModuleServer) readConfig() {
	fileData, err := os.ReadFile("config/" + _self.GetConfigFile())
	if err != nil {
		panic("read config file err")
	}
	_self.config = new(RotaryServerConfig)
	err = json.Unmarshal(fileData, _self.config)
	if err != nil {
		panic("decode config file err")
	}
}

// 注册客户端消息回调
func (_self *RotaryModuleServer) registerClientPacket(clientHandler *ByteConnectionHandler) {
	clientHandler.Register(PacketCommand(0), onRtyMsg)

}

func onRtyMsg(connection Connection, packet *DataPacket) {
	// logger.Info("------------------------raw packet :%v", packet)
	b := _rotaryModuleServer.BaseServer.GetBroker()
	// logger.Info("------++++++++++---cmd:%v -- %v----message  %v  ch", packet.Command(), packet.Message(), newProtoMessage)
	b.Publish(_rotaryHandlerChannel[1], []byte(packet.GetStringData()))
	logger.Info("---data from tcp ---radar---------message  %v  channel:%v", packet.GetStringData(), _rotaryHandlerChannel[1])

}
