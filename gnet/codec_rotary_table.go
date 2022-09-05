package gnet

// const (
// 	ANGLE_SET_TITLE = "SET_angle"
// 	ANGLE_SET_BODY  = "SET_angle_****.**_****.**_****.**"
// 	ANGLE_TITLE     = "RED_angle"
// 	ANGLE_BODY      = "RED_angle_****.**_****.**_****.**"
// 	SPEED_TITLE     = "RED_speed"
// 	SPEED_BODY      = "RED_speed_***.**_***.**_***.**"
// )

// 用了RingBuffer的连接的编解码接口
// 流格式: Length+Data
// 这里把编解码分成了2层
// 第1层:从RingBuffer里取出原始包体数据
// 第2层:对包数据执行实际的编解码操作
type RotaryTableCodec struct {
	// 包头的编码接口,包头长度不能变
	HeaderEncoder func(connection Connection, packet Packet, headerData []byte)
	// 包体的编码接口
	// NOTE:返回值允许返回多个[]byte,如ProtoPacket在编码时,可以分别返回command和proto.Message的序列化[]byte
	// 如果只返回一个[]byte,就需要把command和proto.Message序列化的[]byte再合并成一个[]byte,造成性能损失
	DataEncoder func(connection Connection, packet Packet) [][]byte
	// 包头的解码接口,包头长度不能变
	HeaderDecoder func(connection Connection, headerData []byte)
	// 包体的解码接口
	// DataDecoder func(connection Connection, packetHeader *PacketHeader, packetData []byte) Packet
	DataDecoder func(connection Connection, packetHeader string, packetData []byte) Packet

	// packetHeaderSize PacketHeaderSize
}

func (_self *RotaryTableCodec) CreatePacketHeader(connection Connection, packet Packet, packetData []byte) PacketHeader {
	return NewDefaultPacketHeader(0, 0)
}

func (_self *RotaryTableCodec) PacketHeaderSize() uint32 {
	ROTARY_TITLE := "a"
	// logger.Debug("读取数据头部长度  PacketHeaderSize %v   %v %v ", ROTARY_TITLE, int(unsafe.Sizeof(ROTARY_TITLE)), len(ROTARY_TITLE))
	return uint32(len(ROTARY_TITLE))
}

func (_self *RotaryTableCodec) PacketBodySize(title string) uint32 {
	// logger.Debug("读取数据实体长度  PacketBodySize %v   %v", title, int(unsafe.Sizeof(ANGLE_BODY)))

	// switch title {
	// case ANGLE_TITLE:
	// 	return uint32(len(ANGLE_BODY)) - uint32(len(ANGLE_TITLE))
	// case ANGLE_SET_TITLE:
	// 	return uint32(len(ANGLE_SET_BODY)) - uint32(len(ANGLE_SET_TITLE))
	// case SPEED_TITLE:
	// 	return uint32(len(SPEED_BODY)) - uint32(len(SPEED_TITLE))
	// }
	return 0
}

func (_self *RotaryTableCodec) Encode(connection Connection, packet Packet) []byte {
	if tcpConnection, ok := connection.(*SerialConnection); ok {
		sendBuffer := tcpConnection.sendBuffer
		var encodedData []byte
		if _self.DataEncoder != nil {
			encodedData = _self.DataEncoder(connection, packet)[0]
		} else {
			encodedData = []byte(packet.GetStringData())
		}
		_, err := sendBuffer.Write(encodedData)
		if err != nil {
			logger.Debug("写入字符串数据失败 %v", encodedData)
		}

		return nil
	}

	return []byte{}
}

func (_self *RotaryTableCodec) Decode(connection Connection, data []byte) (newPacket Packet, err error) {
	if tcpConnection, ok := connection.(*SerialConnection); ok {
		recvBuffer := tcpConnection.recvBuffer

		// 先解码包头
		if tcpConnection.curReadPacketHeader == nil {
			packetHeaderSize := int(_self.PacketHeaderSize())
			if recvBuffer.UnReadLength() < packetHeaderSize {
				return
			}
			var packetHeaderData []byte
			readBuffer := recvBuffer.ReadBuffer()
			if len(readBuffer) >= packetHeaderSize {

				packetHeaderData = readBuffer[0:packetHeaderSize]
				// logger.Info("---packetHeaderSize----- %v ", packetHeaderData)
			} else {

				packetHeaderData = tcpConnection.tmpReadPacketHeaderData
				// 先拷贝RingBuffer的尾部
				n := copy(packetHeaderData, readBuffer)
				// 再拷贝RingBuffer的头部
				copy(packetHeaderData[n:], recvBuffer.buffer)
			}
			recvBuffer.SetReaded(packetHeaderSize)

			PLCTitle := string(packetHeaderData)

			// header := tcpConnection.curReadPacketHeader
			bodySize := _self.PacketBodySize(PLCTitle)
			// logger.Info("检查数据包是否读完    recvBuffer.UnReadLength()%v   int(header.Len()):%v  body : %v", PLCTitle, uint32(len(PLCTitle)), bodySize)
			var packetData []byte

			if int(bodySize) <= recvBuffer.Size() {

				if recvBuffer.UnReadLength() < int(bodySize) {
					logger.Info("包体数据还没收完整")
					return
				}
				packetData = recvBuffer.ReadFull(int(bodySize))
			} else {
				// 数据包超出了RingBuffer大小
				// 为什么要处理数据包超出RingBuffer大小的情况?
				// 因为RingBuffer是一种内存换时间的解决方案,对于处理大量连接的应用场景,内存也是要考虑的因素
				// 有一些应用场景,大部分数据包都不大,但是有少量数据包非常大,如果RingBuffer必须设置的比最大数据包还要大,可能消耗过多内存
			}

			if _self.DataDecoder != nil {
				newPacket = _self.DataDecoder(connection, PLCTitle, packetData)
			} else {
				newPacket = NewDataPacket(packetData, "")
				tcpConnection.curReadPacketHeader = nil
				return newPacket, ErrNotSupport
			}
			tcpConnection.curReadPacketHeader = nil

			return newPacket, ErrNotSupport

		}
	}

	return nil, ErrNotSupport
}

func NewRotaryTableCodec(packetHeaderSize PacketHeaderSize) *RotaryTableCodec {
	if packetHeaderSize == nil {
		packetHeaderSize = func() uint32 {
			return uint32(DefaultPacketHeaderSize)
		}
	}
	codec := &RotaryTableCodec{}
	codec.DataDecoder = codec.DecodePacket
	// codec.packetHeaderSize = packetHeaderSize
	return codec
}

func (_self *RotaryTableCodec) DecodePacket(connection Connection, packetHeader string, packetData []byte) Packet {
	decodedPacketData := packetData

	return NewDataPacket(nil, packetHeader+string(decodedPacketData))

}
