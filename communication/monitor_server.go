package communication

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"strconv"
	"sync"
	"time"

	. "desay.com/radar-monitor/gnet"
	. "desay.com/radar-monitor/internal"
	"desay.com/radar-monitor/logger"

	pb "desay.com/radar-monitor/pb"

	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc"
)

var (
	_                  Server = (*MonitorViewServer)(nil)
	_monitorViewServer *MonitorViewServer
	_monitorMsgMap     = &SMap{
		Map: make(map[string][]byte),
	}
	keyFile                = flag.String("key_file", "", "The TLS key file")
	port                   = flag.Int("port", 50051, "The server port")
	_monitorHandlerChannel = []string{"FC2_RECEIVE", "FC2_STATUS"}
	//_count                 = 0
)

type SMap struct {
	sync.RWMutex
	Map map[string][]byte
}

func (l *SMap) readMap(key string) ([]byte, bool) {
	l.RLock()
	value, ok := l.Map[key]
	l.RUnlock()
	return value, ok
}

func (l *SMap) writeMap(key string, value []byte) {
	l.Lock()
	l.Map[key] = value
	l.Unlock()
}

type MonitorViewServer struct {
	BaseServer
	// pb.UnimplementedMonitorViewServer
	mu          sync.Mutex // protects routeNotes
	config      *MonitorServerConfig
	_grpcServer *grpc.Server
	_lastAck    *pb.RadarAckMsgType
}

// 登录服配置
type MonitorServerConfig struct {
	BaseServerConfig
}
type AckResponse struct {
	Code int
	Msg  string
	Data interface{}
}

func (_self *MonitorViewServer) crosConfig(w *http.ResponseWriter) *http.ResponseWriter {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	(*w).Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	(*w).Header().Set("content-type", "application/json")
	return w
}

func Fc2ParamAck(rw http.ResponseWriter, req *http.Request) {
	rw = *_monitorViewServer.crosConfig(&rw)
	_respons := make(map[string]interface{})
	_respons["code"] = 0
	_respons["data"] = _monitorViewServer._lastAck
	_respons["msg"] = ""
	jsonU, _ := json.Marshal(_respons)
	// fmt.Println(string(jsonU))
	rw.Write(jsonU)
}

func ConfigFc2Param(rw http.ResponseWriter, req *http.Request) {
	rw = *_monitorViewServer.crosConfig(&rw)
	// 读取配置内容body
	s, _ := ioutil.ReadAll(req.Body)
	jsonstr := string(s)

	cfg := &pb.RadarParamCfgType{}
	_respons := make(map[string]interface{})
	//构造默认返回值
	ack := &pb.RadarAckMsgType{}
	ack.EnuAckCode = pb.RadarAckMsg_CodeType_ACK_UNDEFINED
	ack.U32AAckData = []uint32{0}
	_respons["code"] = 0
	_respons["data"] = ack
	_respons["msg"] = ""
	// logger.Info(" ConfigFc2Param>>>>>>0000>>>>> >>>>>>>>>>>%v ", jsonstr)
	json.Unmarshal([]byte(jsonstr), cfg)

	// 获取消息中间件

	b := _monitorViewServer.BaseServer.GetBroker()
	if b == nil {
		logger.Info("ConfigFc2Param bridge is nil---------")
		_respons["msg"] = "AckCode_ParamCfgType_ERR_UNKNOWN"
		jsonU, _ := json.Marshal(_respons)
		fmt.Println(string(jsonU))
		rw.Write(jsonU)

		return
	}

	// 往雷达发送消息通道推送配置数据
	// logger.Info(" ConfigFc2Param>>>>>11111>>>>>> %v >>>>>>>>>>> ", cfg)
	data, _ := proto.Marshal(cfg)
	b.Publish(_radarHandlerChannel[0], data)
	// logger.Info(" ConfigFc2Param>>>>>>>>>>> %v >>>>>>>>>>> %X", cfg, data)

	// 订阅雷达ack消息
	key := _radarHandlerChannel[1]

	ch, err := b.Subscribe(key)
	if err != nil {
		log.Fatalf("failed to subscrib: %v", err)
		_respons["msg"] = "AckCode_ParamCfgType_ERR_UNKNOWN"
		jsonU, _ := json.Marshal(_respons)
		fmt.Println(string(jsonU))
		rw.Write(jsonU)
		b.Unsubscribe(key, ch)
		return
	}

	_msg, ok := _monitorMsgMap.readMap(key)
	if !ok {
		_msg = []byte{}
	}
	_payl := b.GetPayLoadAsync(ch, 1500)
	if _payl == nil {
		_respons["code"] = 500
		ack.EnuAckCode = pb.RadarAckMsg_CodeType_ACK_UNDEFINED
		_respons["msg"] = "AckCode_ParamCfgType_ERR_UNKNOWN"
	} else {
		_msg = _payl.([]byte)
		_monitorMsgMap.writeMap(key, _msg)
		proto.Unmarshal(_payl.([]byte), ack)
		_respons["msg"] = ack
	}

	// logger.Info("get byte from channel subscrib------- 	time:%v   subm:%v   _msg %v len :%v", util.GetCurrentMS(), ack, _msg, len(_msg))
	_respons["code"] = 0
	// result, _ = m.MarshalToString(ack)
	_respons["data"] = ack
	jsonU, _ := json.Marshal(_respons)
	// fmt.Println(string(jsonU))
	b.Unsubscribe(key, ch)
	rw.Write(jsonU)
}

func HeartBeatReq(rw http.ResponseWriter, req *http.Request) {
	rw = *_monitorViewServer.crosConfig(&rw)
	_respons := make(map[string]interface{})
	_respons["code"] = 0
	_respons["data"] = ""
	_respons["msg"] = ""
	jsonU, _ := json.Marshal(_respons)
	// fmt.Println(string(jsonU))
	// rw.Write(jsonU)
	rw.Write(jsonU)

}

func RotaryCommand(rw http.ResponseWriter, req *http.Request) {
	rw = *_monitorViewServer.crosConfig(&rw)
	// 读取配置内容body
	s, _ := ioutil.ReadAll(req.Body)
	jsonstr := string(s)
	// logger.Info(" RotaryCommand>>>>>>>111>>>> %v >>>>>>>>>>> ", jsonstr)
	cfg := &pb.RotaryCommand{}
	_respons := make(map[string]interface{})
	//构造默认返回值
	ack := &pb.RotaryCommand{}
	ack.CommandPayload = ""

	_respons["code"] = 0
	_respons["data"] = ack
	_respons["msg"] = ""

	json.Unmarshal([]byte(jsonstr), cfg)

	// logger.Info(" RotaryCommand>>>>>>>>>>> %v >>>>>>>>>>> ", _respons)
	// 获取消息中间件

	b := _monitorViewServer.BaseServer.GetBroker()
	if b == nil {
		logger.Info("RotaryCommand bridge is nil---------")
		_respons["msg"] = "AckCode_ParamCfgType_ERR_UNKNOWN"
		jsonU, _ := json.Marshal(_respons)
		fmt.Println(string(jsonU))
		rw.Write(jsonU)
		return
	}

	// data, _ := proto.Marshal(payload)
	str_cmd := cfg.CommandPayload
	moduleType := cfg.ModuleType
	key := _rotaryHandlerChannel[0]
	if moduleType == 2 {
		key = _rotaryModuleHandlerChannel[0]
	}
	b.Publish(key, []byte(str_cmd))

	// logger.Info(" RotaryCommand>>>>>>>>>>> %v >>>>>>>>>>> %v >>>>>>>%v", cfg, str_cmd, key)

	jsonU, _ := json.Marshal(_respons)
	// fmt.Println(string(jsonU))
	rw.Write(jsonU)
}

func RotaryTableCommand(rw http.ResponseWriter, req *http.Request) {
	rw = *_monitorViewServer.crosConfig(&rw)
	// 读取配置内容body
	s, _ := ioutil.ReadAll(req.Body)
	jsonstr := string(s)
	// logger.Info(" RotaryCommand>>>>>>>111>>>> %v >>>>>>>>>>> ", jsonstr)

	_respons := make(map[string]interface{})
	//构造默认返回值

	_respons["code"] = 0
	_respons["data"] = jsonstr
	_respons["msg"] = ""

	// 获取消息中间件

	b := _monitorViewServer.BaseServer.GetBroker()
	if b == nil {
		logger.Info("RotaryCommand bridge is nil---------")
		_respons["msg"] = "AckCode_ParamCfgType_ERR_UNKNOWN"
		jsonU, _ := json.Marshal(_respons)
		rw.Write(jsonU)
		return
	}

	b.Publish(_rotaryTableHandlerChannel[0], []byte(jsonstr))
	jsonU, _ := json.Marshal(_respons)
	rw.Write(jsonU)
}

func ADCRecvStatistic(rw http.ResponseWriter, req *http.Request) {

	rw = *_monitorViewServer.crosConfig(&rw)
	_respons := make(map[string]interface{})
	_respons["code"] = 500
	_respons["data"] = ""
	_respons["msg"] = ""

	handler := _communicateServer.msgHandler
	// logger.Info("ADCRecvStatistic bridge recv ---111----key:%v  handler. :%v ", key, handler.GetStatistic())
	statistic := handler.GetStatistic()
	_respons["code"] = 0
	_respons["data"] = statistic
	jsonU, _ := json.Marshal(_respons)
	rw.Write(jsonU)

}

func ConfigFileSet(rw http.ResponseWriter, req *http.Request) {
	rw = *_monitorViewServer.crosConfig(&rw)
	query := req.URL.Query()
	names, _ := query["fileName"]
	_respons := make(map[string]interface{})
	_respons["code"] = 500
	_respons["data"] = ""
	_respons["msg"] = ""
	_respons["code"] = 0
	b := _monitorViewServer.BaseServer.GetBroker()
	if b == nil {
		logger.Info("Acknowledged bridge is nil---------")
		_respons["msg"] = "Acknowledged bridge is nil----"
		jsonU, _ := json.Marshal(_respons)
		fmt.Println(string(jsonU))
		rw.Write(jsonU)
		return
	}
	b.Publish(_radarHandlerChannel[3], []byte(names[0]))
	jsonU, _ := json.Marshal(_respons)

	// fmt.Println(string(jsonU))
	rw.Write(jsonU)
	// logger.Info("get byte from channel subscrib------- 	_monitorMsgMap[key]:%v   subm:%v  ", key, jsonU)

}

func RotaryCommandAcknowledge(rw http.ResponseWriter, req *http.Request) {
	rw = *_monitorViewServer.crosConfig(&rw)
	query := req.URL.Query()

	names, ok := query["moduleType"]
	_respons := make(map[string]interface{})
	_respons["code"] = 500
	_respons["data"] = ""
	_respons["msg"] = ""
	if !ok || len(names[0]) < 1 {

		_respons["msg"] = "Url Param 'moduleType' is missing"
		jsonU, _ := json.Marshal(_respons)
		fmt.Println(string(jsonU))
		rw.Write(jsonU)
		return
	}

	ModuleType, err := strconv.Atoi(names[0])
	if err != nil {
		_respons["msg"] = " Param 'moduleType' must be a number"
		jsonU, _ := json.Marshal(_respons)
		fmt.Println(string(jsonU))
		rw.Write(jsonU)
		return
	}

	// logger.Info(" RotaryCommand>>>>>>>111>>>> %v >>>>>>>>>>> %t", names, names)

	b := _monitorViewServer.BaseServer.GetBroker()
	if b == nil {
		logger.Info("Acknowledged bridge is nil---------")
		_respons["msg"] = "Acknowledged bridge is nil----"
		jsonU, _ := json.Marshal(_respons)
		fmt.Println(string(jsonU))
		rw.Write(jsonU)
		return
	}
	// 	// SET_angle_-000.00_0000.00_0000.00 ||33 || _0000.00
	// 	// SET_angle_-000.00_0000.00_0000.00_0000.00

	key := _rotaryHandlerChannel[1]
	key_module := _rotaryModuleHandlerChannel[1]

	// logger.Info("rotaryCommandAcknowledge bridge recv ---111----key%v", key)
	updateTicker := time.NewTicker(123 * time.Millisecond)
	defer func() {
		updateTicker.Stop()
	}()

	ch, err := b.Subscribe(key)

	if err != nil {
		logger.Info("failed to subscrib: %v", err)
		_respons["msg"] = "failed to subscrib"
		jsonU, _ := json.Marshal(_respons)
		fmt.Println(string(jsonU))
		rw.Write(jsonU)
		b.Unsubscribe(key, ch)
		return
	}
	count := 10
	ack := &pb.RotaryCommand{}
	for {
		<-updateTicker.C
		if count == 0 {
			break
		}
		current := key
		if ModuleType == 2 {
			current = key_module
		}
		// ModuleType = 2

		ack.ModuleType = uint32(ModuleType) // 转台

		_msg, ok := _monitorMsgMap.readMap(current)
		if !ok {
			_msg = []byte{}
		}
		_payl := b.GetPayLoadAsync(ch, 100)
		count = count - 1

		if _payl == nil {
			continue
		} else {
			_msg = _payl.([]byte)
			_monitorMsgMap.writeMap(current, _msg)

			ack.CommandPayload = string(_msg)

			_respons["code"] = 0
			_respons["data"] = ack
			jsonU, _ := json.Marshal(_respons)
			fmt.Println(string(jsonU))
			rw.Write(jsonU)
			b.Unsubscribe(key, ch)
			logger.Info("get byte from channel subscrib------- 	_monitorMsgMap[key]:%v   subm:%v  ", ack)
			return
		}

	}

	_respons["code"] = 0
	_respons["data"] = ack
	jsonU, _ := json.Marshal(_respons)
	b.Unsubscribe(key, ch)
	// fmt.Println(string(jsonU))
	rw.Write(jsonU)
}

func RotaryTableCommandAcknowledge(rw http.ResponseWriter, req *http.Request) {
	rw = *_monitorViewServer.crosConfig(&rw)

	_respons := make(map[string]interface{})
	_respons["code"] = 500
	_respons["data"] = ""
	_respons["msg"] = ""
	ack := &pb.RotaryCommand{}
	ack.CommandPayload = "{\"X\":0,\"Y\":0}"
	b := _monitorViewServer.BaseServer.GetBroker()
	if b == nil {
		logger.Info("Acknowledged bridge is nil---------")
		_respons["msg"] = "Acknowledged bridge is nil----"
		jsonU, _ := json.Marshal(_respons)
		rw.Write(jsonU)
		return
	}

	b.Publish(_rotaryTableHandlerChannel[0], []byte("POS"))
	key := _rotaryTableHandlerChannel[1]

	updateTicker := time.NewTicker(123 * time.Millisecond)
	defer func() {
		updateTicker.Stop()
	}()

	ch, err := b.Subscribe(key)

	if err != nil {
		logger.Info("failed to subscrib: %v", err)
		_respons["msg"] = "failed to subscrib"
		jsonU, _ := json.Marshal(_respons)
		rw.Write(jsonU)
		b.Unsubscribe(key, ch)
		return
	}
	count := 20
	// ack := &pb.RotaryCommand{}
	for {
		<-updateTicker.C
		if count == 0 {
			break
		}
		current := key

		_msg, ok := _monitorMsgMap.readMap(current)
		if !ok {
			_msg = []byte{}
		}
		_payl := b.GetPayLoadAsync(ch, 100)
		count = count - 1

		if _payl == nil {
			continue
		} else {
			_msg = _payl.([]byte)
			_monitorMsgMap.writeMap(current, _msg)
			ack.CommandPayload = string(_msg)
			_respons["code"] = 0
			_respons["data"] = ack
			jsonU, _ := json.Marshal(_respons)
			rw.Write(jsonU)
			b.Unsubscribe(key, ch)
			logger.Info("get byte from channel subscrib------- 	_monitorMsgMap[key]:%v   subm:%v  ", ack)
			return
		}

	}

	_respons["code"] = 0
	_respons["data"] = ack
	jsonU, _ := json.Marshal(_respons)
	b.Unsubscribe(key, ch)
	rw.Write(jsonU)
}

func newServer() *MonitorViewServer {
	s := &MonitorViewServer{}
	return s
}

func registerReceiveChannel() {
	s := _monitorViewServer
	b := s.BaseServer.GetBroker()
	if b == nil {
		return
	}
	for _, key := range _monitorHandlerChannel {
		ch, err := b.Subscribe(key)
		if err != nil {
			log.Fatalf("failed to subscrib: %v", err)
		}

		go func(ch <-chan interface{}, key string) {
			for {
				_msg, ok := _monitorMsgMap.readMap(key)
				if !ok {
					_msg = []byte{}
				}
				_payl := b.GetPayLoad(ch)

				_msg = _payl.([]byte)
				_monitorMsgMap.writeMap(key, _msg)
				// copy(_msg, _payl.([]byte))
				ack := &pb.RadarAckMsgType{}
				proto.Unmarshal(_payl.([]byte), ack)
				logger.Info("get byte from channel subscrib------- 	_monitorMsgMap[key]:%v   subm:%v   total:%v ", _msg, ack)
				if key == _monitorHandlerChannel[1] {
					s._lastAck = ack

				}
			}
		}(ch, key)
	}
}

// 初始化
func (_self *MonitorViewServer) Init(ctx context.Context, configFile string) bool {
	_monitorViewServer = _self

	_self.BaseServer.AddInitHook(registerReceiveChannel)

	if !_self.BaseServer.Init(ctx, configFile) {
		return false
	}

	clientCodec := NewProtoCodec(nil)
	clientHandler := NewDefaultConnectionHandler(clientCodec)
	_self.registerClientPacket(clientHandler)
	flag.Parse()
	_self._lastAck = &pb.RadarAckMsgType{}
	// var opts []grpc.ServerOption
	// if *tls {
	// 	if *certFile == "" {
	// 		*certFile = data.Path("x509/server_cert.pem")
	// 	}
	// 	if *keyFile == "" {
	// 		*keyFile = data.Path("x509/server_key.pem")
	// 	}
	// 	creds, err := credentials.NewServerTLSFromFile(*certFile, *keyFile)
	// 	if err != nil {
	// 		log.Fatalf("Failed to generate credentials %v", err)
	// 	}
	// 	opts = []grpc.ServerOption{grpc.Creds(creds)}
	// }
	// _self._grpcServer = grpc.NewServer(opts...)
	// pb.RegisterMonitorViewServer(_self._grpcServer, newServer())

	// lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", *port))
	// if err != nil {
	// 	log.Fatalf("failed to listen: %v", err)
	// }
	// _self._grpcServer.Serve(lis)
	addr := fmt.Sprintf("127.0.0.1:%d", *port)
	logger.Info("ListenAndServe   0000   total:%v ", addr)
	http.HandleFunc("/ConfigFc2Param", ConfigFc2Param)
	http.HandleFunc("/RotaryCommand", RotaryCommand)
	http.HandleFunc("/HeartBeatReq", HeartBeatReq)
	http.HandleFunc("/RotaryCommandAcknowledge", RotaryCommandAcknowledge)
	http.HandleFunc("/ADCRecvStatistic", ADCRecvStatistic)
	http.HandleFunc("/ConfigFileSet", ConfigFileSet)
	http.HandleFunc("/Fc2ParamAck", Fc2ParamAck)
	http.HandleFunc("/RotaryTableCommandAcknowledge", RotaryTableCommandAcknowledge)
	http.HandleFunc("/RotaryTableCommand", RotaryTableCommand)
	go _monitorViewServer.runView()

	http.ListenAndServe(addr, nil)
	logger.Info("ListenAndServe   111   total:%v ", addr)

	return true
}

func (_self *MonitorViewServer) runView() {
	// logger.Info("ListenAndServe   111   ---------run ---- ")
	// cmd := exec.Command("D:/project/pangoo-radar-monitor-view/out/pangoo-radar-monitor-view-win32-x64/pangoo-radar-monitor-view.exe")
	cmd := exec.Command("pangoo-radar-monitor-view.exe")
	err := cmd.Run()
	if err != nil {
		logger.Error("Error: %v", err)
	}
}

func (_self *MonitorViewServer) Run(ctx context.Context) {

}

// 退出
func (_self *MonitorViewServer) Exit() {
	_self.BaseServer.Exit()
	logger.Info("MonitorViewServer.Exit")

}

// 读取配置文件
func (_self *MonitorViewServer) readConfig() {
}

// 注册客户端消息回调
func (_self *MonitorViewServer) registerClientPacket(clientHandler *DefaultConnectionHandler) {

}
