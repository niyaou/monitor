package gnet

import (
	"context"
	"io"
	"net"
	"sync"
	"time"

	"github.com/jacobsa/go-serial/serial"
	"google.golang.org/protobuf/proto"
)

type SerialConnection struct {
	baseConnection
	// conn net.Conn
	conn io.ReadWriteCloser

	// 防止执行多次关闭操作
	closeOnce sync.Once
	// 关闭回调
	onClose func(connection Connection)
	// 最近收到完整数据包的时间(时间戳:秒)
	lastRecvPacketTick uint32
	// 发包缓存chan
	sendPacketCache chan Packet
	// 发包RingBuffer
	sendBuffer *RingBuffer
	// 收包RingBuffer
	recvBuffer *RingBuffer
	// 解码时,用到的一些临时变量
	tmpReadPacketHeader     PacketHeader
	tmpReadPacketHeaderData []byte
	curReadPacketHeader     PacketHeader
	tmpReadUnDecodeData     []byte
	tmpReadedIndex          int
	dataCount               int
	readLoopTotal           uint64 //总共读取的字节数
	codecLen                uint32 //proto包的长度
	readedHeadLen           uint32 //已经读取了的tcp头长度
	readedDataLen           uint32 //已经读取了的TCP包DATA长度
	isFixedHead             bool   //true: 接收固定长度head的tcp包，用于protobuf，  false:接收字符串
	SerialName              string //串口名
	SerialPort              uint   //串口码率
}

func NewSerialConnector(SerialName string, SerialPort uint, config *ConnectionConfig, codec Codec, handler ConnectionHandler, isFixedHead bool) *SerialConnection {
	// logger.Info("---NewTcpConnector---++++++-")
	if config.MaxPacketSize == 0 {
		config.MaxPacketSize = MaxPacketDataSize
	}
	if config.MaxPacketSize > MaxPacketDataSize {
		config.MaxPacketSize = MaxPacketDataSize
	}
	newConnection := createSerialConnection(SerialName, SerialPort, config, codec, handler)
	newConnection.isConnector = true
	newConnection.isFixedHead = isFixedHead

	return newConnection
}

func createSerialConnection(SerialName string, SerialPort uint, config *ConnectionConfig, codec Codec, handler ConnectionHandler) *SerialConnection {
	// logger.Info("---createTcpConnection---+++++++-")
	newConnection := &SerialConnection{
		baseConnection: baseConnection{
			addr:         "",
			connectionId: newConnectionId(),
			config:       config,
			codec:        codec,
			handler:      handler,
		},
		sendPacketCache: make(chan Packet, config.SendPacketCacheCap),
	}
	newConnection.SerialName = SerialName
	newConnection.SerialPort = SerialPort
	newConnection.tmpReadPacketHeader = codec.CreatePacketHeader(newConnection, nil, nil)
	// logger.Info("---createTcpConnection---+++++++- %v %v", newConnection, newConnection.tmpReadPacketHeader)
	return newConnection
}

func NewSerialConnectionAccept(conn net.Conn, config *ConnectionConfig, codec Codec, handler ConnectionHandler) *SerialConnection {
	// logger.Info("---NewTcpConnectionAccept---+++++++-")
	if config.MaxPacketSize == 0 {
		config.MaxPacketSize = MaxPacketDataSize
	}
	if config.MaxPacketSize > MaxPacketDataSize {
		config.MaxPacketSize = MaxPacketDataSize
	}
	return &SerialConnection{
		baseConnection: baseConnection{
			connectionId: newConnectionId(),
			isConnector:  false,
			config:       config,
			codec:        codec,
			handler:      handler,
		},
		sendPacketCache: make(chan Packet, config.SendPacketCacheCap),
		conn:            conn,
	}
}

// 连接
func (_self *SerialConnection) Connect(address string) bool {
	options := serial.OpenOptions{
		PortName:        _self.SerialName,
		BaudRate:        _self.SerialPort,
		DataBits:        8,
		StopBits:        1,
		MinimumReadSize: 4,
	}

	// conn, err := net.DialTimeout("tcp", address, time.Second)
	logger.Info("serial  Open %v %v", _self.SerialName, _self.SerialPort)
	conn, err := serial.Open(options)
	if err != nil {
		_self.isConnected = false

		return false
	}
	_self.conn = conn
	_self.isConnected = true
	if _self.handler != nil {
		_self.handler.OnConnected(_self, true)
	}
	return true
}

// 开启读写协程
func (_self *SerialConnection) Start(ctx context.Context, netMgrWg *sync.WaitGroup, onClose func(connection Connection)) {
	// logger.Info("-++++++++--TcpConnection-----Start---- %v", _self.GetConnectionId())
	// 开启收包协程
	netMgrWg.Add(1)
	go func() {
		defer func() {
			netMgrWg.Done()
			if err := recover(); err != nil {
				logger.Error("read fatal %v: %v", _self.GetConnectionId(), err.(error))
				LogStack()
			}
		}()
		_self.readLoop()
		logger.Info("readLoop   %v %v", _self.SerialName, _self.SerialPort)

		_self.Close()
	}()

	// 开启发包协程
	netMgrWg.Add(1)
	go func(ctx context.Context) {
		defer func() {
			netMgrWg.Done()
			if err := recover(); err != nil {
				logger.Error("write fatal %v: %v", _self.GetConnectionId(), err.(error))
				LogStack()
			}
		}()
		_self.writeLoop(ctx)
		_self.Close()
	}(ctx)
}

// 收包过程
func (_self *SerialConnection) readLoop() {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("readLoop fatal %v: %v", _self.GetConnectionId(), err.(error))
			LogStack()
		}
	}()
	_self.recvBuffer = _self.createRecvBuffer()
	_self.tmpReadPacketHeaderData = make([]byte, _self.codec.PacketHeaderSize())
	for _self.isConnected {
		// 可写入的连续buffer
		writeBuffer := _self.recvBuffer.WriteBuffer()
		readBuffer := _self.recvBuffer.ReadBuffer()

		n, err := _self.conn.Read(writeBuffer)
		if err != nil {
			if err != io.EOF {
				logger.Debug("readLoop %v err:%v", _self.GetConnectionId(), err.Error())
			}
			break
		}
		_self.recvBuffer.SetWrited(n)
		logger.Info(" count:%v   unread length %v ", n, _self.recvBuffer.UnReadLength())

		for _self.isConnected {
			newPacket, _ := _self.codec.Decode(_self, readBuffer)
			if newPacket == nil {
				break
			}
			_self.lastRecvPacketTick = GetCurrentTimeStamp()
			if _self.handler != nil {
				_self.dataCount++
				_self.handler.OnRecvPacket(_self, newPacket)
			}
		}
	}

}

// 发包过程
func (_self *SerialConnection) writeLoop(ctx context.Context) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("writeLoop fatal %v: %v", _self.GetConnectionId(), err.(error))
			//LogStack()
		}
		logger.Debug("writeLoop end %v", _self.GetConnectionId())
	}()

	// logger.Debug("writeLoop begin %v", _self.GetConnectionId())
	// 收包超时计时,用于检测掉线
	recvTimeoutTimer := time.NewTimer(time.Second * time.Duration(_self.config.RecvTimeout))
	defer recvTimeoutTimer.Stop()
	// 心跳包计时
	heartBeatTimer := time.NewTimer(time.Second * time.Duration(_self.config.HeartBeatInterval))
	defer heartBeatTimer.Stop()
	_self.sendBuffer = _self.createSendBuffer()
	for _self.isConnected {
		var delaySendDecodePacketData []byte
		select {
		case packet := <-_self.sendPacketCache:
			if packet == nil {
				logger.Debug("packet==nil %v", _self.GetConnectionId())
				return
			}
			// 数据包编码
			// Encode里面会把编码后的数据直接写入sendBuffer
			delaySendDecodePacketData = _self.codec.Encode(_self, packet)
			// logger.Debug("-----数据包编码 -----结果------- %v   %v   %v", packet.Command(), packet.GetStringData(), delaySendDecodePacketData)
			if len(delaySendDecodePacketData) > 0 {
				// Encode里面写不完的数据延后处理
				logger.Debug("%v sendBuffer is full delaySize:%v", _self.GetConnectionId(), len(delaySendDecodePacketData))
				break
			}
			packetCount := len(_self.sendPacketCache)
			// 还有其他数据包在chan里,就进行批量合并
			if packetCount > 0 {
				for i := 0; i < packetCount; i++ {
					// 这里不会阻塞
					newPacket, ok := <-_self.sendPacketCache
					if !ok {
						logger.Debug("newPacket==nil %v", _self.GetConnectionId())
						return
					}
					// 数据包编码
					delaySendDecodePacketData = _self.codec.Encode(_self, newPacket)
					if len(delaySendDecodePacketData) > 0 {
						logger.Debug("%v sendBuffer is full delaySize:%v", _self.GetConnectionId(), len(delaySendDecodePacketData))
						break
					}
				}
			}

		case <-recvTimeoutTimer.C:
			if _self.config.RecvTimeout > 0 {
				nextTimeoutTime := _self.config.RecvTimeout + _self.lastRecvPacketTick - GetCurrentTimeStamp()
				if nextTimeoutTime > 0 {
					recvTimeoutTimer.Reset(time.Second * time.Duration(nextTimeoutTime))
				} else {
					// 指定时间内,一直未读取到数据包,则认为该连接掉线了,可能处于"假死"状态了
					// 需要主动关闭该连接,防止连接"泄漏"
					logger.Debug("recv timeout %v", _self.GetConnectionId())
					return
				}
			}

		case <-heartBeatTimer.C:
			if _self.isConnector && _self.config.HeartBeatInterval > 0 && _self.handler != nil {
				if heartBeatPacket := _self.handler.CreateHeartBeatPacket(_self); heartBeatPacket != nil {
					delaySendDecodePacketData = _self.codec.Encode(_self, heartBeatPacket)
					heartBeatTimer.Reset(time.Second * time.Duration(_self.config.HeartBeatInterval))
				}
			}

		case <-ctx.Done():
			// 收到外部的关闭通知
			logger.Debug("recv closeNotify %v", _self.GetConnectionId())
			return
		}

		if _self.sendBuffer.UnReadLength() > 0 {
			// logger.Info("ring buffer 有内容，需要写入")
			// 可读数据有可能分别存在数组的尾部和头部,所以需要循环发送,有可能需要发送多次
			for _self.isConnected && _self.sendBuffer.UnReadLength() > 0 {
				// if _self.config.WriteTimeout > 0 {
				// 	setTimeoutErr := _self.conn.SetWriteDeadline(time.Now().Add(time.Duration(_self.config.WriteTimeout) * time.Second))
				// 	if setTimeoutErr != nil {
				// 		logger.Debug("%v setTimeoutErr:%v", _self.GetConnectionId(), setTimeoutErr.Error())
				// 		return
				// 	}
				// }
				readBuffer := _self.sendBuffer.ReadBuffer()
				writeCount, err := _self.conn.Write(readBuffer)
				if err != nil {

					logger.Debug("%v write Err:%v", _self.GetConnectionId(), err.Error())
					return
				}
				// logger.Info("ring buffer 有内容，需要写入 writeCount %v  ", writeCount)
				_self.sendBuffer.SetReaded(writeCount)
				if len(delaySendDecodePacketData) > 0 {
					writedLen, _ := _self.sendBuffer.Write(delaySendDecodePacketData)
					// 这里不一定能全部写完
					if writedLen < len(delaySendDecodePacketData) {
						delaySendDecodePacketData = delaySendDecodePacketData[writedLen:]
						logger.Debug("%v write delaybuffer :%v", _self.GetConnectionId(), writedLen)
					} else {
						delaySendDecodePacketData = nil
					}
				}
			}
		}
	}
}

// 关闭
func (_self *SerialConnection) Close() {
	_self.closeOnce.Do(func() {
		_self.isConnected = false
		if _self.conn != nil {
			_self.conn.Close()
			logger.Debug("close %v", _self.GetConnectionId())
			//_self.conn = nil
		}
		if _self.handler != nil {
			_self.handler.OnDisconnected(_self)
		}
		if _self.onClose != nil {
			_self.onClose(_self)
		}
	})
}

// 异步发送proto包
// NOTE:调用Send(command,message)之后,不要再对message进行读写!
func (_self *SerialConnection) Send(command PacketCommand, message proto.Message) bool {
	if !_self.isConnected {
		return false
	}
	packet := NewProtoPacket(command, message)
	// logger.Info(" ---------after encode command %v   message %v", packet.Command(), packet.Message())

	// NOTE:当sendPacketCache满时,这里会阻塞
	_self.sendPacketCache <- packet
	return true
}

// 异步发送数据
// NOTE:调用SendPacket(packet)之后,不要再对packet进行读写!
func (_self *SerialConnection) SendPacket(packet Packet) bool {
	if !_self.isConnected {
		return false
	}
	// NOTE:当sendPacketCache满时,这里会阻塞
	// logger.Info(" ---------after encode command %v   message %v", packet.Command(), packet.GetStringData())
	_self.sendPacketCache <- packet
	return true
}

// 超时发包,超时未发送则丢弃,适用于某些允许丢弃的数据包
// 可以防止某些"不重要的"数据包造成chan阻塞,比如游戏项目常见的聊天广播
func (_self *SerialConnection) TrySendPacket(packet Packet, timeout time.Duration) bool {
	if timeout == 0 {
		// 非阻塞方式写chan
		select {
		case _self.sendPacketCache <- packet:
			return true
		default:
			return false
		}
	}
	sendTimeout := time.After(timeout)
	for {
		select {
		case _self.sendPacketCache <- packet:
			return true
		case <-sendTimeout:
			return false
		}
	}
	return false
}

// 创建用于批量发包的RingBuffer
func (_self *SerialConnection) createSendBuffer() *RingBuffer {
	ringBufferSize := _self.config.SendBufferSize
	if ringBufferSize == 0 {
		if _self.config.MaxPacketSize > 0 {
			ringBufferSize = _self.config.MaxPacketSize * 2
		} else {
			ringBufferSize = 65535
		}
	}
	return NewRingBuffer(int(ringBufferSize))
}

// 创建用于批量收包的RingBuffer
func (_self *SerialConnection) createRecvBuffer() *RingBuffer {
	ringBufferSize := _self.config.RecvBufferSize
	if ringBufferSize == 0 {
		if _self.config.MaxPacketSize > 0 {
			ringBufferSize = _self.config.MaxPacketSize * 2
		} else {
			ringBufferSize = 65535
		}
	}
	return NewRingBuffer(int(ringBufferSize))
}

// 发包RingBuffer
func (_self *SerialConnection) GetSendBuffer() *RingBuffer {
	return _self.sendBuffer
}

// 收包RingBuffer
func (_self *SerialConnection) GetRecvBuffer() *RingBuffer {
	return _self.recvBuffer
}

// LocalAddr returns the local network address.
func (_self *SerialConnection) LocalAddr() net.Addr {
	if _self.conn == nil {
		return nil
	}
	return nil
}

// RemoteAddr returns the remote network address.
func (_self *SerialConnection) RemoteAddr() string {

	// if _self.conn == nil {
	// 	return nil
	// }
	return _self.baseConnection.addr
}

func (_self *SerialConnection) GetSendPacketChanLen() int {
	return len(_self.sendPacketCache)
}
