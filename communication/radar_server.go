package communication

import (
	"context"
	"encoding/json"
	"os"
	"strings"

	. "desay.com/radar-monitor/gnet"
	. "desay.com/radar-monitor/handler"
	. "desay.com/radar-monitor/internal"
	"desay.com/radar-monitor/logger"
	pb "desay.com/radar-monitor/pb"
	"google.golang.org/protobuf/proto"
)

var (
	_                    Server = (*CommunicateServer)(nil)
	_radarHandlerChannel        = []string{"FC2_SEND", "FC2_ACK", "ROTARY_SEND", "CONFIG_FILE", "FC2_STATUS"} //雷达服务处理的频道key
	_radarMsgMap                = &SMap{
		Map: make(map[string][]byte),
	} //每个频道对应的map数据
	_communicateServer *CommunicateServer
)

type CommunicateServer struct {
	BaseServer
	config     *CommunicateServerConfig
	conn       Connection
	netMgr     *NetMgr
	msgHandler *ADCHandler
	dataCount  int
}

type CommunicateServerConfig struct {
	BaseServerConfig
}

func radarSendChannel() {
	s := _communicateServer
	b := s.BaseServer.GetBroker()
	if b == nil {
		return
	}
	for _, key := range _radarHandlerChannel {
		ch, err := b.Subscribe(key)
		if err != nil {
			logger.Error("failed to subscrib: %v", err)
		} else {
			// logger.Info("subscrib channel %v  successed     ", key)
		}

		go func(ch <-chan interface{}, key string) {
			for {
				_msg, ok := _radarMsgMap.readMap(key)
				if !ok {
					_msg = []byte{}
				}
				_payl := b.GetPayLoad(ch)

				_msg = _payl.([]byte)
				_radarMsgMap.writeMap(key, _msg)

				//如果是发送，则构造消息体发送给雷达
				switch key {
				case _radarHandlerChannel[0]:
					point := &pb.RadarParamCfgType{}
					proto.Unmarshal(_msg, point)
					// logger.Info("get from web-------newProtoMessage ---------  pb:%v", point)
					newProtoMessage := &pb.RadarHc_PayloadType{}
					newProtoMessage.U16ProjId = 1
					newProtoMessage.EnuMsgGenre = pb.MsgGenreType_REQUEST
					newProtoMessage.EnuMsgModel = pb.MsgModelType_RADAR_PARAM_CFG
					// data, _ := proto.Marshal(point)
					newProtoMessage.PMsgStream = _msg
					_communicateServer.msgHandler.SetFrameSaveCount(point.U16AcqNrFrames)
					if point.U16Cmd == 4 || point.U16Cmd == 3 {
						_communicateServer.msgHandler.SetFrameSaveCount(200)
					}
					if point.U16Cmd == 65535 {
						_communicateServer.msgHandler.SetCeaseFlag(true)
						logger.Info("send-------newProtoMessage --------- SetCeaseFlag")
					}
					// _parse_data, _ := proto.Marshal(newProtoMessage)

					if s.conn != nil {

						logger.Info("send-------newProtoMessage --------- connect %s  addr %s", s.conn.IsConnected(), s.conn.RemoteAddr())
						if s.conn.IsConnected() {
							s.conn.Send(PacketCommand(1), newProtoMessage)
						} else {
							isConnect := s.netMgr.StartConnector(s.BaseServer.GetContext(), s.conn.RemoteAddr(), s.conn)
							if isConnect {
								s.conn.Send(PacketCommand(1), newProtoMessage)
							}
						}
					}
				case _radarHandlerChannel[2]:
					//获取设置转台的参数，用于设置存储的ADC文件名
					_str_cmd := string(_msg)
					if find := strings.Contains(_str_cmd, "SET_angle_"); find {
						_str_cmd = strings.Replace(_str_cmd, "SET_angle_", "", 1)
						_communicateServer.msgHandler.SetSaveFileName(_str_cmd)
						logger.Info("获取设置转台的参数，用于设置存储的ADC文件名 %v  ", _str_cmd)
					}
				case _radarHandlerChannel[3]:
					_str_cmd := string(_msg)
					_communicateServer.msgHandler.SetConfigFileName(_str_cmd)
					logger.Info("获取设置雷达的参数，-----用于设置存储的ADC文件名 %v  ", _str_cmd)
				default:
					logger.Info("receive-------newProtoMessage ---------default %v  ", key)
				}

			}
		}(ch, key)
	}
}

func onRadarConnected(connection Connection, success bool) {
	// if success {
	// 	_communicateServer.conn = connection
	// }
}

// 初始化
func (_self *CommunicateServer) Init(ctx context.Context, configFile string) bool {
	_communicateServer = _self

	_self.BaseServer.AddInitHook(radarSendChannel)

	if !_self.BaseServer.Init(ctx, configFile) {
		return false
	}
	_self.readConfig()

	_self.netMgr = GetNetMgr()
	_self.dataCount = 0

	clientCodec := NewProtoCodec(nil)
	// clientCodec := NewDefaultCodec()
	clientHandler := NewDefaultConnectionHandler(clientCodec)
	clientHandler.SetOnConnectedFunc(onRadarConnected)
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
	b := _communicateServer.BaseServer.GetBroker()
	_self.msgHandler = NewADCHandler(ctx, b)
	_communicateServer.conn = _self.netMgr.NewConnectorCustom(ctx, _self.config.ClientListenAddr, &_self.config.ClientConnConfig, clientCodec,
		clientHandler, nil, func(config *ConnectionConfig, addr string, codec Codec, handler ConnectionHandler, isFixedHead bool) Connection {
			return NewTcpConnector(config, addr, codec, handler, true)
		}, true)
	if _communicateServer.conn == nil {
		panic("connect failed")
	}

	return true
}

// 运行
func (_self *CommunicateServer) Run(ctx context.Context) {
	_self.BaseServer.Run(ctx)
}

// 退出
func (_self *CommunicateServer) Exit() {
	_self.BaseServer.Exit()
}

// 读取配置文件
func (_self *CommunicateServer) readConfig() {
	fileData, err := os.ReadFile("config/" + _self.GetConfigFile())
	if err != nil {
		panic("read config file err")
	}
	_self.config = new(CommunicateServerConfig)
	err = json.Unmarshal(fileData, _self.config)
	if err != nil {
		panic("decode config file err")
	}
}

// 注册客户端消息回调
func (_self *CommunicateServer) registerClientPacket(clientHandler *DefaultConnectionHandler) {
	clientHandler.Register(PacketCommand(1), onEthanMsg, func() proto.Message { return &pb.RadarHc_PayloadType{} })
}

// 客户端心跳回复
func onHeartBeatReq(connection Connection, packet *ProtoPacket) {
}

func onEthanMsg(connection Connection, packet *ProtoPacket) {
	/**CommunicateServer
	上位机：
	u16ProjId   = 1;    FC2
	enuMsgGenre = 0;    REQUEST
	enuMsgModel = 1;    RADAR_PARAM_CFG
	pMsgStream  = strem;    RADAR_PARAM_CFG消息

	ECU：
	u16ProjId   = 1;   /* FC2
	enuMsgGenre = 1;   /* RESPONSE
	enuMsgModel = 1;   /* RADAR_PARAM_CFG
	pMsgStream  = strem;   /* RADAR_PARAM_CFG应答消息
	*/
	packet.Command()

	newProtoMessage := packet.Message().(*pb.RadarHc_PayloadType)
	b := _communicateServer.BaseServer.GetBroker()
	if newProtoMessage.U16ProjId == 1 {
		switch newProtoMessage.EnuMsgGenre {
		case pb.MsgGenreType_REQUEST:
			b.Publish(_monitorHandlerChannel[0], newProtoMessage.PMsgStream)
		case pb.MsgGenreType_RESPONSE:
			switch newProtoMessage.EnuMsgModel {
			case pb.MsgModelType_RADAR_PARAM_CFG:
				ack := &pb.RadarAckMsgType{}
				proto.Unmarshal(newProtoMessage.PMsgStream, ack)
				logger.Info("------radar---------message------>MsgModelType_RADAR_PARAM_CFG---->  %v", ack)
				b.Publish(_radarHandlerChannel[1], newProtoMessage.PMsgStream)
				b.Publish(_radarHandlerChannel[4], newProtoMessage.PMsgStream)
			case pb.MsgModelType_RADAR_ADC_ACQ_DATA:
				_communicateServer.dataCount++
				_communicateServer.msgHandler.Consume(newProtoMessage.PMsgStream)
			default:
				logger.Info("------radar---------message------>adc---->  dataCount:%v", _communicateServer.dataCount)
			}

		default:
			logger.Info("------no handler---------===================    :")
		}
	}

}
