package handler

import (
	"bufio"
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"os"

	. "desay.com/radar-monitor/internal"
	"desay.com/radar-monitor/logger"
	pb "desay.com/radar-monitor/pb"
	. "desay.com/radar-monitor/util"
	"github.com/golang/protobuf/proto"
	"github.com/hedzr/go-ringbuf/v2"
	"github.com/hedzr/go-ringbuf/v2/mpmc"
)

type Handler interface {
	Consume(p []byte) error
	Dismiss() error
	Start() error
	LoopTackle() error
}

type ADCHandler struct {
	cacheRing       mpmc.RingBuffer[[]byte]
	dealChan        chan []byte
	runningFlag     bool
	ctx             context.Context
	datFile         *os.File
	wDatFileIo      *bufio.Writer
	datFileBig      *os.File
	wDatFileBigIo   *bufio.Writer
	fileDataCount   int
	dataHandling    int
	frameHeadCount  int
	chirpHeadCount  int
	chirpDataCount  int
	chirpTailCount  int
	totalBytesCount uint32
	totalFrameLen   uint32
	frameSaveCount  uint32 //每次文件保存的帧数
	saveFrameCount  uint32 //已经保存的帧数
	broker          *BrokerClient
	channelRunning  bool
	ceaseSaveFlag   bool   //停止信号，收到停止信号后，结束当前文件
	saveFileName    string //存储的文件名
	configFileName  string //当前配置文件名
}

func NewADCHandler(ctx context.Context, broker *BrokerClient) *ADCHandler {
	// dealChan := make(chan []byte)
	dealChan := make(chan []byte, 1024)
	return &ADCHandler{
		cacheRing:      ringbuf.New[[]byte](12000),
		runningFlag:    false,
		fileDataCount:  0,
		dataHandling:   0,
		dealChan:       dealChan,
		broker:         broker,
		channelRunning: false,
		frameSaveCount: 1,
		saveFrameCount: 0,
		ceaseSaveFlag:  false,
		configFileName: "",
	}
}

func errChk(err error) {
	if err != nil {
		logger.Error("failed to enqueue : %v", err)
	}
}

func (_self *ADCHandler) parseMessageAndWrite(_msg proto.Message, fileBufIo *bufio.Writer, isLittle bool) {
	tail := proto.MessageReflect(_msg)
	tailDesc := tail.Descriptor()
	fieldDs := tailDesc.Fields()
	for i := 0; i < fieldDs.Len(); i++ {
		fieldD := fieldDs.Get(i)
		val := tail.Get(fieldD) // 获取字段值
		if fieldD.IsList() {
			u32arr := []uint32{}
			for j := 0; j < val.List().Len(); j++ {
				u32arr = append(u32arr, uint32(val.List().Get(j).Uint()))
			}
			bytesBuffer := bytes.NewBuffer([]byte{})
			if isLittle {
				binary.Write(bytesBuffer, binary.LittleEndian, u32arr)
			} else {
				binary.Write(bytesBuffer, binary.BigEndian, u32arr)
			}
			byteData := bytesBuffer.Bytes()
			fileBufIo.Write(byteData)

			_self.totalBytesCount = _self.totalBytesCount + uint32(len(byteData))

		} else {
			if fieldD.Kind().String() == "uint64" {
				_byteDate := []byte{0, 0, 0, 0, 0, 0, 0, 0}
				if isLittle {
					binary.LittleEndian.PutUint64(_byteDate, uint64(val.Uint()))
				} else {
					binary.BigEndian.PutUint64(_byteDate, uint64(val.Uint()))
				}
				fileBufIo.Write(_byteDate)

				_self.totalBytesCount = _self.totalBytesCount + uint32(len(_byteDate))

			} else {
				byteDate := []byte{0, 0, 0, 0}
				if isLittle {
					binary.LittleEndian.PutUint32(byteDate, uint32(val.Uint()))
				} else {
					binary.BigEndian.PutUint32(byteDate, uint32(val.Uint()))
				}
				fileBufIo.Write(byteDate)

				_self.totalBytesCount = _self.totalBytesCount + uint32(len(byteDate))

			}
		}
	}
}

func (_self *ADCHandler) GetStatistic() []uint32 {
	statistic := []uint32{uint32(_self.chirpTailCount), _self.totalFrameLen, _self.totalBytesCount}
	return statistic
}

func (_self *ADCHandler) SetFrameSaveCount(count uint32) {
	_self.frameSaveCount = count
}

func (_self *ADCHandler) SetCeaseFlag(flag bool) {
	_self.ceaseSaveFlag = flag
}

func (_self *ADCHandler) SetSaveFileName(turrentParams string) {
	_self.saveFileName = turrentParams
}

func (_self *ADCHandler) SetConfigFileName(configFileName string) {
	_self.configFileName = configFileName
}

func (_self *ADCHandler) LoopTackle() error {
	// 读取数据，写入文件
	if _self.channelRunning {
		return nil
	}
	_self.channelRunning = true
	go func() {

		for item := range _self.dealChan {
			adc := &pb.AdcAcq_MsgType{}
			err := proto.Unmarshal(item, adc)
			if err != nil {
				logger.Info("Consume----read error-%v ", err)
			}
			switch adc.EnuAdcAcqDataType {
			case pb.AdcAcq_DataType_ADC_ACQ_TYPE_FRAME_HEAD:
				if _self.saveFrameCount == 0 {
					_self.datFileBig, _ = os.OpenFile(fmt.Sprintf("AdcAcqData/%s%s_%s.dat", _self.configFileName, GetCurrentTimeString(), _self.saveFileName), os.O_WRONLY|os.O_CREATE, os.ModePerm)
					_self.wDatFileBigIo = bufio.NewWriter(_self.datFileBig)
					_self.fileDataCount = 0
					_self.dataHandling = 0
					_self.frameHeadCount = 0
					_self.chirpHeadCount = 0
					_self.chirpDataCount = 0
					_self.chirpTailCount = 0
					_self.totalBytesCount = 0
					_self.totalFrameLen = 0
					logger.Info("----------》>>>>>>>>>>>>>>>>>open file ============== %v    dealing:%v  ", fmt.Sprintf("AdcAcqData/%s_param_%s_big.dat", GetCurrentTimeString(), _self.saveFileName))
				}

				frameHead := &pb.AdcAcq_FrameHeadType{}
				proto.Unmarshal(adc.PDataStream, frameHead)

				_self.parseMessageAndWrite(frameHead, _self.wDatFileBigIo, false)
				_self.totalFrameLen = _self.totalFrameLen + frameHead.U32FrameLen
				_self.dataHandling++
				_self.frameHeadCount++
				logger.Info("----------》decode ============_FRAME_HEAD")
				// logger.Info("----------》AdcAcq_DataType_ADC_ACQ_TYPE_FRAME_HEAD ============== data receive:%v    dealing:%v   frameHeadCount %v  chirpHeadCount:%v   chirpDataCount:%v   chirpTailCount:%v _self.totalBytesCount:%v",
				// 	_self.fileDataCount, _self.dataHandling, _self.frameHeadCount, _self.chirpHeadCount, _self.chirpDataCount, _self.chirpTailCount, _self.totalBytesCount)
			case pb.AdcAcq_DataType_ADC_ACQ_TYPE_CHIRP_HEAD:
				chirpHead := &pb.AdcAcq_ChirpHeadType{}
				proto.Unmarshal(adc.PDataStream, chirpHead)
				_self.parseMessageAndWrite(chirpHead, _self.wDatFileBigIo, false)
				_self.dataHandling++
				_self.chirpHeadCount++
				// logger.Info("----------》decode =============_CHIRP_HEAD")
			case pb.AdcAcq_DataType_ADC_ACQ_TYPE_CHIRP_DATA:
				chirpData := &pb.AdcAcq_ChirpChSampleDataType{}
				proto.Unmarshal(adc.PDataStream, chirpData)

				_self.parseMessageAndWrite(chirpData, _self.wDatFileBigIo, false)
				_self.dataHandling++
				_self.chirpDataCount++
				// logger.Info("----------》decode =============_CHIRP_data")
			case pb.AdcAcq_DataType_ADC_ACQ_TYPE_CHIRP_TAIL:
				chirpTail := &pb.AdcAcq_FrameTailType{}
				proto.Unmarshal(adc.PDataStream, chirpTail)
				_self.parseMessageAndWrite(chirpTail, _self.wDatFileBigIo, false)
				_self.dataHandling++
				_self.chirpTailCount++
				_self.saveFrameCount++
				logger.Info("----------》decode =============_CHIRP_TAIL")
				// _self.publicMsg(statistic)
				// logger.Info("-------------------Publish -------------------")
				if _self.saveFrameCount == _self.frameSaveCount || _self.ceaseSaveFlag {
					_self.Dismiss()
					_self.saveFrameCount = 0
					_self.ceaseSaveFlag = false
					// logger.Info("----------dismiss %v ===================", _self.saveFrameCount)
					logger.Info("----------》dequeue end  ======_self.saveFrameCount %v == _self.frameSaveCount  %v ============= data receive:%v    dealing:%v   frameHeadCount %v  chirpHeadCount:%v   chirpDataCount:%v   chirpTailCount:%v totalBytesCount:%v  totalFrameLen:%v",
						_self.saveFrameCount, _self.frameSaveCount, _self.fileDataCount, _self.dataHandling, _self.frameHeadCount, _self.chirpHeadCount, _self.chirpDataCount, _self.chirpTailCount, _self.totalBytesCount, _self.totalFrameLen)

				}

			default:
				logger.Info("------------》dequeue ok  》: %v   msg:%v   byte: % 02X  \n", adc.EnuAdcAcqDataType, []uint32{1}, adc.PDataStream)
			}

		}

	}()

	return nil
}

func (_self *ADCHandler) Start() error {
	if !_self.runningFlag {
		_self.runningFlag = true
		_self.LoopTackle()
	}

	return nil
}

func (_self *ADCHandler) Consume(p []byte) error {
	_self.dealChan <- p
	_self.fileDataCount++
	err := _self.Start()
	errChk(err)
	return nil
}

func (_self *ADCHandler) Dismiss() error {
	_self.wDatFileBigIo.Flush()
	_self.datFileBig.Close()
	_self.runningFlag = false
	return nil
}
