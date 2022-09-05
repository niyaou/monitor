package main

import (
	"encoding/binary"
	"net"
	"time"

	"desay.com/radar-monitor/logger"
	pb "desay.com/radar-monitor/pb"
	"github.com/golang/protobuf/proto"
)

func send(point uint32, conn net.Conn) {

	var msg *pb.RadarHc_PayloadType = &pb.RadarHc_PayloadType{}
	msg.U16ProjId = 1
	msg.EnuMsgGenre = pb.MsgGenreType_RESPONSE
	msg.EnuMsgModel = pb.MsgModelType_RADAR_PARAM_CFG
	var payload = &pb.TEF82Xx_FrameOptionalParamType{}
	payload.U32SeqInterval = point
	payload.U8SweepRstCtrl = point
	payload.U8CafcPllLPFSel = point

	data, err := proto.Marshal(payload)
	if err != nil {
		logger.Error("proto encode err:%v cmd:%v", err, data)
		return
	}
	msg.PMsgStream = data
	// logger.Info(" msg %v", msg)

	src, _ := proto.Marshal(msg)
	// srcLen := len(src)
	rangeBuff := []byte{byte(len(src)), 0, 0, 0}
	rangeBuff = append(rangeBuff, src...)
	// logger.Info("src:%v  srcLen:%v  rangeBuff:%v len:%v", src, srcLen, rangeBuff, len(rangeBuff), srcLen)
	// logger.Info("payload %v   rangeBuff:%v rangeBuff:%v", payload, rangeBuff, len(rangeBuff))
	binary.LittleEndian.PutUint32(rangeBuff, uint32(len(src)))
	conn.Write(rangeBuff)

	testLitBuff := []byte{0, 0, 0, 0}
	testBigBuff := []byte{0, 0, 0, 0}
	binary.LittleEndian.PutUint32(testLitBuff, uint32(256*256*256))
	binary.BigEndian.PutUint32(testBigBuff, uint32(256*256*256))
	// encodedStr := hex.EncodeToString(src)
	logger.Info("binary encode testLitBuff:%v    testBigBuff %v ", testLitBuff, testBigBuff)
	// data1, _ := hex.DecodeString(encodedStr)
	// for _, v := range data1 {
	// 	fmt.Printf("%X \n", v)
	// }
}

func main() {

	// connect to server
	// conn, _ := net.Dial("tcp", "127.0.0.1:10558")
	conn, _ := net.Dial("tcp", "127.0.0.1:10003")
	// logger.Info("conn  %v", conn)
	// fps := []uint32{}
	// fps := []uint32{1, 2, 3, 4, 5, 6, 7, 8, 9, 1, 2, 4, 1, 2, 3, 4, 5, 6, 7, 8, 9, 1, 2, 4}

	updateTicker := time.NewTicker(100 * time.Millisecond)
	defer func() {
		updateTicker.Stop()
	}()

	for i := 1; i <= 1; i++ {
		<-updateTicker.C
		send(uint32(i), conn)
	}
	// fps := []uint32{1, 2, 3, 4, 5}
	// for _, v := range fps {
	// println(fmt.Sprintf("i: %v \n", i))
	// println(fmt.Sprintf("v: %v \n", v))
	// send(v, conn)
	// }

	// v := 0
	// for {

	// 	t := <-updateTicker.C
	// 	fmt.Println("当前时间为:", t)
	// 	if v >= 100 {
	// 		break
	// 	}
	// 	send(uint32(v), conn)
	// 	v++
	// }
}
