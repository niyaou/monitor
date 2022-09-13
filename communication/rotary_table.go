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
	_                          Server = (*RotaryTableServer)(nil)
	_rotaryTableHandlerChannel        = []string{"ROTARY_TABLE_SEND", "ROTARY_TABLE_ACK"}
	_rotaryTableMsgMap                = make(map[string][]byte)
	// singleton
	_rotaryTableServer *RotaryTableServer
)

type RotaryTableServer struct {
	BaseServer
	config *RotaryServerConfig
	conn   Connection
	netMgr *NetMgr
}

type RotaryTableServerConfig struct {
	BaseServerConfig
}

func rotaryTableSendChannel() {
	s := _rotaryTableServer
	b := s.BaseServer.GetBroker()
	if b == nil {
		return
	}
	for _, key := range _rotaryTableHandlerChannel {

		ch, err := b.Subscribe(key)
		if err != nil {
			logger.Error("failed to subscrib: %v", err)
		} else {
		}

		go func(ch <-chan interface{}, key string) {
			for {
				_msg, ok := _rotaryTableMsgMap[key]
				if !ok {
					_msg = []byte{}
				}
				_payl := b.GetPayLoad(ch)
				_msg = _payl.([]byte)
				_str_cmd := string(_msg)

				switch key {
				case _rotaryTableHandlerChannel[0]:
					logger.Info("send-------to channel  to device ---------% 02X str:%v  ", _msg, _str_cmd)
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
				case _rotaryHandlerChannel[1]:
					break
				}

			}
		}(ch, key)
	}
}

func onTableConnected(connection Connection, success bool) {
	if success {
		_rotaryTableServer.conn = connection

	}

}

// 初始化
func (_self *RotaryTableServer) Init(ctx context.Context, configFile string) bool {
	_rotaryTableServer = _self

	_self.BaseServer.AddInitHook(rotaryTableSendChannel)

	if !_self.BaseServer.Init(ctx, configFile) {
		return false
	}
	_self.readConfig()

	_self.netMgr = GetNetMgr()

	// clientCodec := NewProtoCodec(nil)
	clientCodec := NewRotaryTableCodec(nil)
	// clientCodec := NewRotaryModuleCodec(nil)

	clientHandler := NewByteConnectionHandler(clientCodec)

	clientHandler.SetOnConnectedFunc(onTableConnected)
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

	_rotaryTableServer.conn = _self.netMgr.NewConnectorCustom(ctx, _self.config.ClientListenAddr, &_self.config.ClientConnConfig, clientCodec,
		clientHandler, nil, func(config *ConnectionConfig, addr string, codec Codec, handler ConnectionHandler, isFixedHead bool) Connection {
			return NewSerialConnector(_self.config.SerialName, _self.config.SerialPort, config, codec, handler, false)
		}, false)
	if _rotaryTableServer.conn == nil {
		panic("connect failed")
	}

	b := _self.BaseServer.GetBroker()
	b.Publish(_rotaryTableHandlerChannel[0], []byte("POS"))
	return true
}

// 运行
func (_self *RotaryTableServer) Run(ctx context.Context) {
	_self.BaseServer.Run(ctx)
	// logger.Info("LoginServer.Run")
}

// 退出
func (_self *RotaryTableServer) Exit() {
	_self.BaseServer.Exit()
}

// 读取配置文件
func (_self *RotaryTableServer) readConfig() {
	fileData, err := os.ReadFile("config/" + _self.GetConfigFile())
	if err != nil {
		panic("read config file err")
	}
	_self.config = new(RotaryServerConfig)
	err = json.Unmarshal(fileData, _self.config)
	if err != nil {
		panic("decode config file err")
	}
	// logger.Debug("%v", _self.config)
	// _self.BaseServer.GetServerInfo().ServerId = _self.config.ServerId
	// _self.BaseServer.GetServerInfo().ServerType = "login"
	// _self.BaseServer.GetServerInfo().ClientListenAddr = _self.config.ClientListenAddr
}

// 注册客户端消息回调
func (_self *RotaryTableServer) registerClientPacket(clientHandler *ByteConnectionHandler) {
	clientHandler.Register(PacketCommand(0), onSPMsg)

}

func onSPMsg(connection Connection, packet *DataPacket) {
	// logger.Info("------------------------raw packet :%v", packet)
	b := _rotaryTableServer.BaseServer.GetBroker()
	// logger.Info("------++++++++++---cmd:%v -- %v----message  %v  ch", packet.Command(), packet.Message(), newProtoMessage)
	b.Publish(_rotaryHandlerChannel[1], []byte(packet.GetStringData()))
	logger.Info("---data from tcp ---radar---------message  %v  channel:%v", packet.GetStringData(), _rotaryHandlerChannel[1])

}
