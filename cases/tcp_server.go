package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"math"
	"net"
	"strconv"
	"strings"
	"time"

	"desay.com/radar-monitor/logger"
	pb "desay.com/radar-monitor/pb"
	"github.com/golang/protobuf/proto"
)

func sendString(conn net.Conn) {
	ANGLE_BODY := "RED_angle_****.**_****.**_****.**"
	conn.Write([]byte(ANGLE_BODY))
	logger.Info("Client 建立连接，发来的数据 ---------- string:%v n %v ", ANGLE_BODY, []byte(ANGLE_BODY))
}

func processRadar(conn net.Conn) {
	defer conn.Close() // 关闭连接
	az := int64(-5)
	el := int64(0)
	polar := int64(0)

	az_des := int64(-5)
	el_des := int64(0)
	polar_des := int64(0)

	for {
		reader := bufio.NewReader(conn)
		var buf [10240]byte
		n, err := reader.Read(buf[:])
		if err != nil {
			fmt.Println("read from client failed, err: ", err)
			break
		} else {
			recvStr := string(buf[:n])
			logger.Info("收到Client端发来的数据：", recvStr)
			if len(recvStr) > 8 && strings.Compare(recvStr[:9], "SET_angle") == 0 {
				//设置参数
				//         10 11	17 19 	25 27	33
				// SET_angle_+000.00_+002.00_+003.00
				_polar_des, _ := strconv.ParseFloat(recvStr[11:17], 32)
				polar_des = int64(math.Floor(_polar_des))
				if strings.Compare(recvStr[10:11], "-") == 0 {
					logger.Info(">>>>>>>--11---- polar_des:%v %v", polar_des, recvStr[11:17])
					polar_des = polar_des - polar_des*2
					logger.Info(">>>>>>>------ polar_des:%v %v", polar_des, recvStr[11:17])
				}
				_el_des, _ := strconv.ParseFloat(recvStr[19:25], 32)
				el_des = int64(math.Floor(_el_des))
				if strings.Compare(recvStr[18:19], "-") == 0 {
					logger.Info(">>>>>>>--11---- el_des:%v %v", el_des, recvStr[19:25])
					el_des = el_des - el_des*2
					logger.Info(">>>>>>>------ el_des:%v %v", el_des, recvStr[19:25])
				}
				_az_des, _ := strconv.ParseFloat(recvStr[27:33], 32)
				az_des = int64(math.Floor(_az_des))
				if strings.Compare(recvStr[26:27], "-") == 0 {
					logger.Info(">>>>>>>---11--- az_des:%v %v", az_des, recvStr[27:33])
					az_des = az_des - az_des*2
					logger.Info(">>>>>>>------ az_des:%v %v", az_des, recvStr[27:33])
				}
				logger.Info(">>>>>>>------ %v  %v  %v", polar_des, el_des, az_des)

			}

			if strings.Compare(recvStr, "RED_angle") == 0 {
				logger.Info(">>>>>>>------ polar %v    polar_des  %v     el  %v    el_des  %v    az  %v    az_des %v", polar, polar_des, el, el_des, az, az_des)

				str := "RED_angle_"
				if polar == polar_des && el == el_des && az == az_des {

					str = "SET_angle_"
				}
				if polar < 0 {
					str = str + "-0%02d.00_"
				} else {
					str = str + "00%02d.00_"
				}

				if el < 0 {
					str = str + "-0%02d.00_"
				} else {
					str = str + "00%02d.00_"
				}
				if az < 0 {
					str = str + "-0%02d.00"
				} else {
					str = str + "00%02d.00"
				}
				// logger.Info(">>>>>>>------ string:%v %v", str, fmt.Sprintf(str, int(math.Abs(float64(polar))), int(math.Abs(float64(el))), int(math.Abs(float64(az)))))
				conn.Write([]byte(fmt.Sprintf(str, int(math.Abs(float64(polar))), int(math.Abs(float64(el))), int(math.Abs(float64(az))))))

				if polar == polar_des {
				} else if polar_des > polar {
					polar = polar + 1
				} else if polar_des < polar {
					polar = polar - 1
				}

				if el == el_des {

				} else if el_des > el {
					el = el + 1
				} else if el_des < el {
					el = el - 1

				}

				if az == az_des {

				} else if az_des > az {
					az = az + 1
				} else if az_des < az {
					az = az - 1
				}

			}
		}
	}
}

func processRadarMOdual(conn net.Conn) {
	defer conn.Close() // 关闭连接
	anglen := 0
	for {
		reader := bufio.NewReader(conn)
		var buf [1024]byte
		n, err := reader.Read(buf[:])
		if err != nil {
			fmt.Println("read from client failed, err: ", err)
			break
		} else {
			recvStr := string(buf[:n])

			fmt.Println("收到Client端发来的数据：", recvStr, strings.Compare(recvStr, "RED_angle") == 0)
			if strings.Compare(recvStr, "RED_angle") == 0 {
				logger.Info(">>>>>>>------ string:%v", fmt.Sprintf("RED_angle_0000.0%d_0000.0%d_0000.0%d_0000.0%d_0000.0%d", anglen, anglen, anglen, anglen, anglen))
				conn.Write([]byte(fmt.Sprintf("RED_angle_0000.0%d_0000.0%d_0000.0%d_0000.0%d_0000.0%d", anglen, anglen, anglen, anglen, anglen)))
				anglen = anglen + 1
				if anglen > 9 {
					anglen = 0
				}
			}
		}
	}
}

func pushMessage(conn net.Conn, newProtoMessage *pb.RadarHc_PayloadType, payload proto.Message, count int) int {
	//发送count++12
	buf, _ := proto.Marshal(payload)
	newProtoMessage.PMsgStream = buf

	_msg, _ := proto.Marshal(newProtoMessage)
	lens := []byte{0, 0, 0, 0}
	binary.LittleEndian.PutUint32(lens, uint32(len(_msg)))
	_msg = append(lens, _msg...)
	// logger.Info("构造newProtoMessage %v   len %v  _msg:%v", lens, binary.LittleEndian.Uint32(lens), _msg)
	conn.Write(_msg) // 发送数据
	// _size := 12
	// if len(buf) < 12 {
	// 	_size = len(buf)
	// }
	count = count + 1
	// logger.Info("发送数据raw %v    buf % 02x     ", payload.(*pb.AdcAcq_MsgType).EnuAdcAcqDataType, buf[:_size])
	return count

}

func processAdc(conn net.Conn) {
	defer conn.Close() // 关闭连接
	for {
		reader := bufio.NewReader(conn)
		var buf [1024]byte

		_, err := reader.Read(buf[:]) // 读取数据

		if err != nil {
			fmt.Println("read from client failed, err: ", err)
			break
		} else {
			plen := binary.LittleEndian.Uint32(buf[:4])
			// logger.Info("++++++收到Client端发来的数据：%v  len:%v n %v ", buf[:4], plen, n)

			newProtoMessage := &pb.RadarHc_PayloadType{}
			proto.Unmarshal(buf[4:plen+5], newProtoMessage)

			cfg := &pb.RadarParamCfgType{}
			proto.Unmarshal(newProtoMessage.PMsgStream, cfg)
			framesCount := cfg.U16AcqNrFrames
			if cfg.U16Cmd == 4 || cfg.U16Cmd == 3 {
				framesCount = 200
			}

			byteCount := uint32(0)
			count := 0
			frameHeadCount := 0
			chirpHeadCount := 0
			chirpDataCount := 0
			chirpTailCount := 0
			chirpTotalLen := uint32(0)
			for i := uint32(0); i < framesCount; i++ {
				newProtoMessage.U16ProjId = 1
				newProtoMessage.EnuMsgGenre = pb.MsgGenreType_RESPONSE
				newProtoMessage.EnuMsgModel = pb.MsgModelType_RADAR_ADC_ACQ_DATA
				//构造ADC数据
				adc := &pb.AdcAcq_MsgType{}
				adc.EnuAdcAcqDataType = pb.AdcAcq_DataType_ADC_ACQ_TYPE_FRAME_HEAD

				frameHead := &pb.AdcAcq_FrameHeadType{}
				frameHead.U32FrameStart = 0xAAAAAAAA
				frameHead.U32FrameSeq = 1
				frameHead.U32FrameLen = 1
				frameHead.U64LocalTime = 202208051711
				frameHead.U64TimeStamp = 152534777412
				frameHead.U32FrameInfo = 552
				frameHead.U32ChirpInfo = 125

				chirpHead := &pb.AdcAcq_ChirpHeadType{}
				chirpHead.U32HeadInfo = 652

				chirpData := &pb.AdcAcq_ChirpChSampleDataType{}
				// chirpData.U16ChirpData_Ch1 = []uint32{1, 2, 6, 54, 7, 55, 47, 98, 12, 45, 12, 50}
				chirpData.U32AChSampleData = []uint32{}
				for i := 0; i < 512; i++ {
					// chirpData.U16ChirpData_Ch1[i] = uint32(i)
					chirpData.U32AChSampleData = append(chirpData.U32AChSampleData, uint32(0x0f04))
				}

				chirpTail := &pb.AdcAcq_FrameTailType{}
				chirpTail.U32EndFlag1 = 0xffffffff
				chirpTail.U32EndFlag2 = 0xffffffff

				frameLen := 0
				buf, _ := proto.Marshal(frameHead)
				frameLen += len(buf)
				buf, _ = proto.Marshal(chirpHead)
				frameLen += len(buf)
				buf, _ = proto.Marshal(chirpData)
				chirpDataArr := 8
				frameLen += len(buf) * chirpDataArr
				buf, _ = proto.Marshal(chirpTail)
				frameLen += len(buf)

				frameHead.U32FrameLen = uint32(frameLen)
				chirpTotalLen = chirpTotalLen + uint32(frameHead.U32FrameLen)
				adc.PDataStream, _ = proto.Marshal(frameHead)
				frameHeadCount++

				byteCount = byteCount + uint32(len(adc.PDataStream))
				count = pushMessage(conn, newProtoMessage, adc, count)
				// logger.Info("开始发送数据chirpData----- %v    frameHeadCount %v  chirpHeadCount:%v   chirpDataCount:%v   chirpTailCount:%v byteCount:%v",
				// 	count, frameHeadCount, chirpHeadCount, chirpDataCount, chirpTailCount, byteCount)

				for i := 0; i < 512; i++ {
					adc.EnuAdcAcqDataType = pb.AdcAcq_DataType_ADC_ACQ_TYPE_CHIRP_HEAD
					adc.PDataStream, _ = proto.Marshal(chirpHead)
					byteCount = byteCount + uint32(len(adc.PDataStream))
					count = pushMessage(conn, newProtoMessage, adc, count)
					chirpHeadCount++

					ticker := time.NewTicker(1000 * time.Nanosecond)
					dataCount := 0
					for {
						<-ticker.C
						adc.EnuAdcAcqDataType = pb.AdcAcq_DataType_ADC_ACQ_TYPE_CHIRP_DATA
						adc.PDataStream, _ = proto.Marshal(chirpData)
						byteCount = byteCount + uint32(len(adc.PDataStream))
						count = pushMessage(conn, newProtoMessage, adc, count)
						chirpDataCount++
						dataCount++

						if dataCount >= chirpDataArr {
							break
						}
					}
					ticker.Stop()
				}

				adc.EnuAdcAcqDataType = pb.AdcAcq_DataType_ADC_ACQ_TYPE_CHIRP_TAIL
				adc.PDataStream, _ = proto.Marshal(chirpTail)

				byteCount = byteCount + uint32(len(adc.PDataStream))
				count = pushMessage(conn, newProtoMessage, adc, count)
				chirpTailCount++
			}
			logger.Info("开始发送数据----- %v    frameHeadCount %v  chirpHeadCount:%v   chirpDataCount:%v   chirpTailCount:%v  byteCount:%v  chirpTotalLen:%v",
				count, frameHeadCount, chirpHeadCount, chirpDataCount, chirpTailCount, byteCount, chirpTotalLen)
		}
	}
}

// TCP Server端测试
// 处理函数
func process(conn net.Conn) {
	defer conn.Close() // 关闭连接
	for {
		reader := bufio.NewReader(conn)
		var buf [1024]byte

		// _, err1 := reader.Read(buf[0:4]) // 读取数据

		// if err1 != nil {
		// 	fmt.Println("read from client failed, err: ", err1)
		// 	break
		// }

		n, err := reader.Read(buf[:]) // 读取数据

		if err != nil {
			fmt.Println("read from client failed, err: ", err)
			break
		} else {
			// logger.Info("---------收到Client端发来的数据：%v  ", buf)
			newProtoMessage := &pb.RadarHc_PayloadType{}
			plen := binary.LittleEndian.Uint32(buf[:4])
			proto.Unmarshal(buf[4:plen+5], newProtoMessage)
			// logger.Info("---------收到Client端发来的数据1111：%v  ", newProtoMessage)
			cfg := &pb.RadarParamCfgType{}
			proto.Unmarshal(newProtoMessage.PMsgStream, cfg)

			logger.Info("++++++收到Client端发来的数据：%v  len:%v n %v ", buf[:4], plen, n)
			logger.Info("收到Client端发来的数据：cfg %v   ", cfg)

			newProtoMessage.U16ProjId = 1
			newProtoMessage.EnuMsgGenre = pb.MsgGenreType_RESPONSE
			newProtoMessage.EnuMsgModel = pb.MsgModelType_RADAR_PARAM_CFG

			ack := &pb.RadarAckMsgType{}
			ack.EnuAckCode = pb.RadarAckMsg_CodeType_ACK_UNDEFINED
			ack.U32AAckData = []uint32{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
			logger.Info("构造ACK %v   ", ack)

			buf, err := proto.Marshal(ack)
			if err != nil {
				logger.Error("proto encode err:%v cmd:%v", err, buf)
				return
			}
			newProtoMessage.PMsgStream = buf
			logger.Info("构造newProtoMessage %v   buf %v", newProtoMessage, buf)
			_msg, _ := proto.Marshal(newProtoMessage)
			lens := []byte{0, 0, 0, 0}
			binary.LittleEndian.PutUint32(lens, uint32(len(_msg)))
			_msg = append(lens, _msg...)
			// conn.Write(buf[:n]) // 发送数据
			conn.Write(_msg) // 发送数据
			logger.Info("发送数据raw %v    byte %v   ", newProtoMessage, _msg)
		}
		// recvStr := string(buf[:n])
		// fmt.Println("收到Client端发来的数据：", recvStr)
		// conn.Write([]byte(recvStr)) // 发送数据
	}
}

func main() {

	listen, err := net.Listen("tcp", "0.0.0.0:10003")
	// listen, err := net.Listen("tcp", "0.0.0.0:55562")
	// listen2, err := net.Listen("tcp", "0.0.0.0:3000")
	// listen3, err := net.Listen("tcp", "0.0.0.0:3001")
	if err != nil {
		fmt.Println("Listen() failed, err: ", err)
		return
	}

	// chirpTail := &pb.AdcAcq_FrameTailType{}
	// chirpTail.U32EndFlag1 = 0xffff
	// chirpTail.U32EndFlag2 = 0xffff
	// data := *(*[]byte)(unsafe.Pointer(chirpTail))
	// mash, _ := proto.Marshal(chirpTail)
	// logger.Info("发送数据       byte % 02X   proto:% 02X  msg: %v ", data, mash, chirpTail)
	// Len := unsafe.Sizeof(chirpTail)
	// chirpData := &pb.AdcAcq_ChirpChSampleDataType{}
	// chirpData.U32AChSampleData = []uint32{1, 2, 6, 54, 7, 55, 47, 98, 12, 45, 12, 50}
	// tail := proto.MessageReflect(chirpData)

	// frameHead := &pb.AdcAcq_FrameTailType{}
	// frameHead.U32EndFlag1 = 0xffffffff
	// frameHead.U32EndFlag2 = 0xffffffff

	// frameHead := &pb.AdcAcq_FrameHeadType{}

	// frameHead.U32FrameStart = 0xAAAAAAAA
	// frameHead.U32FrameSeq = 1
	// frameHead.U32FrameLen = 1
	// frameHead.U64LocalTime = 202208051711
	// frameHead.U64TimeStamp = 152534777412
	// frameHead.U32FrameInfo = 552
	// frameHead.U32ChirpInfo = 125
	// tail := proto.MessageReflect(frameHead)
	//
	// tailDesc := tail.Descriptor()

	// logger.Info("获取字段值   message :%v    ", tailDesc.FullName())
	// fieldDs := tailDesc.Fields()
	// for i := 0; i < fieldDs.Len(); i++ {
	// 	fieldD := fieldDs.Get(i)

	// 	val := tail.Get(fieldD) // 获取字段值
	// 	byteDate := []byte{0, 0, 0, 0}
	// 	binary.BigEndian.PutUint32(byteDate, uint32(val.Uint()))
	// 	fmt.Printf("分析当前proto字段  name: %s   fullname: %s   kind: %v   isList :%v   val:%s  byte % 02X   list % 02X  \n", fieldD.Name(), fieldD.FullName(), fieldD.Kind(), fieldD.IsList(), val, val, byteDate)
	// 	// logger.Info("获取字段值   message :%v    byte % 02X   name: %v  ", frameHead, val, fieldD.Name())
	// 	// 	if fieldD.IsList() {
	// 	// 		for j := 0; j < val.List().Len(); j++ {
	// 	// 			logger.Info("数组值  %v", val.List().Get(j))
	// 	// 		}
	// 	// 	}

	// }

	for {
		conn, err := listen.Accept() // 监听客户端的连接请求
		// conn2, err := listen2.Accept() // 监听客户端的连接请求
		// conn3, err := listen3.Accept() // 监听客户端的连接请求
		if err != nil {
			fmt.Println("Accept() failed, err: ", err)
			continue
		}
		go process(conn) // 启动一个goroutine来处理客户端的连接请求
		// go processAdc(conn)
		// go processRadar(conn)
		// go processRadarMOdual(conn)
		// go sendString(conn) // 启动一个goroutine来处理客户端的连接请求
	}

}
