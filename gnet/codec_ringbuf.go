package gnet

// 用了RingBuffer的连接的编解码接口
// 流格式: Length+Data
// 这里把编解码分成了2层
// 第1层:从RingBuffer里取出原始包体数据
// 第2层:对包数据执行实际的编解码操作
type RingBufferCodec struct {
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
	DataDecoder func(connection Connection, packetHeader PacketHeader, packetData []byte) Packet

	packetHeaderSize PacketHeaderSize
}
type PacketHeaderSize func() uint32

func (_self *RingBufferCodec) CreatePacketHeader(connection Connection, packet Packet, packetData []byte) PacketHeader {
	return NewDefaultPacketHeader(0, 0)
}

func (_self *RingBufferCodec) PacketHeaderSize() uint32 {
	// logger.Debug("  ------ring buff-------  PacketHeaderSize------")
	return uint32(DefaultPacketHeaderSize)
}

func (_self *RingBufferCodec) Encode(connection Connection, packet Packet) []byte {
	// 优化思路:编码后的数据直接写入RingBuffer.sendBuffer,可以减少一些内存分配
	if tcpConnection, ok := connection.(*TcpConnection); ok {
		packetHeaderSize := int(_self.packetHeaderSize())
		sendBuffer := tcpConnection.sendBuffer
		var encodedData [][]byte
		if _self.DataEncoder != nil {
			// 编码接口可能把消息分解成了几段字节流数组,如消息头和消息体
			// 如果只是返回一个[]byte结果的话,那么编码接口还需要把消息头和消息体进行合并,从而多一次内存分配和拷贝
			encodedData = _self.DataEncoder(connection, packet)
			// logger.Debug("数据包编码   encodedData %v   %v len %v", 33, encodedData, len(encodedData))
		} else {
			// 支持在应用层做数据包的序列化和编码
			encodedData = [][]byte{packet.GetStreamData()}
		}
		encodedDataLen := 0
		for i, data := range encodedData {
			if i == 0 {
				continue
			}
			// logger.Debug("---获取数据包长度   Encode %v   %v", encodedDataLen, data)
			encodedDataLen += len(data)
		}
		// packetHeader := NewDefaultPacketHeader(uint32(encodedDataLen), 0)
		// packetHeaderSize := int(unsafe.Sizeof(DefaultPacketHeader{}))
		packetHeader := NewDefaultPacketHeader(uint32(encodedDataLen), 0)
		writeBuffer := sendBuffer.WriteBuffer()
		if packetHeaderSize == DefaultPacketHeaderSize && len(writeBuffer) >= packetHeaderSize {
			// logger.Debug("有足够的连续空间可写,则直接写入RingBuffer里    %v  %v %v ", encodedData, len(encodedData), len(writeBuffer))
			// 	// 有足够的连续空间可写,则直接写入RingBuffer里
			// 	// 省掉了一次内存分配操作: make([]byte, PacketHeaderSize)
			packetHeader.WriteTo(writeBuffer)
			// if _self.HeaderEncoder != nil {
			// 	_self.HeaderEncoder(connection, packet, writeBuffer[0:packetHeaderSize])
			// 	logger.Debug("_self.HeaderEncoder != nil 写入header %v   ", writeBuffer)
			// }
			sendBuffer.SetWrited(packetHeaderSize)
			// logger.Debug("SetWrited encodedData %v  encodedDataLen  %v", encodedData, encodedDataLen)
		} else {
			logger.Debug("没有足够的连续空间可写,则能写多少写多少,有可能一部分写入尾部,一部分写入头部 %v   %v", 44, encodedData)
			// // 没有足够的连续空间可写,则能写多少写多少,有可能一部分写入尾部,一部分写入头部
			// packetHeaderData := make([]byte, packetHeaderSize)
			// packetHeader.WriteTo(packetHeaderData)
			// if _self.HeaderEncoder != nil {
			// 	_self.HeaderEncoder(connection, packet, packetHeaderData)
			// }
			// writedHeaderLen, _ := sendBuffer.Write(packetHeaderData)
			// if writedHeaderLen < packetHeaderSize {
			// 	// 写不下的包头数据和包体数据,返回给TcpConnection延后处理
			// 	// 合理的设置发包缓存,一般不会运行到这里
			// 	remainData := make([]byte, packetHeaderSize-writedHeaderLen+encodedDataLen)
			// 	// 没写完的header数据
			// 	n := copy(remainData, packetHeaderData[writedHeaderLen:])
			// 	// 编码后的包体数据
			// 	for _, data := range encodedData {
			// 		n += copy(remainData[n:], data)
			// 	}
			// 	return remainData
			// }
		}
		writedDataLen := 0
		// logger.Debug("range encodedData-------- %v   %v  %v", encodedData, packetHeader, packetHeader.Len())
		for i, data := range encodedData {
			// logger.Debug("write to buffer--- %v  %v", data, i)
			if i == 0 {
				continue
			}
			writed, _ := sendBuffer.Write(data)
			writedDataLen += writed
			// 	if writed < len(data) {
			// 		// 写不下的包体数据,返回给TcpConnection延后处理
			// 		remainData := make([]byte, encodedDataLen-writedDataLen)
			// 		n := copy(remainData, data[writed:])
			// 		for j := i + 1; j < len(encodedData); j++ {
			// 			n += copy(remainData[n:], encodedData[j])
			// 		}
			// 		return remainData
			// 	}
		}
		return nil
	}

	return packet.GetStreamData()
}

func (_self *RingBufferCodec) Decode(connection Connection, data []byte) (newPacket Packet, err error) {
	if tcpConnection, ok := connection.(*TcpConnection); ok {
		// TcpConnection用了RingBuffer,解码时,尽可能的不产生copy

		if _self.DataDecoder != nil {
			// 包体的解码接口
			newPacket = _self.DataDecoder(connection, nil, data)
			// logger.Info("包体的解码接口 +++++++++++ %v", newPacket.Message())
		}
		tcpConnection.curReadPacketHeader = nil
		return newPacket, ErrNotSupport

	}
	// logger.Info("解码接口 return nil ---->>>>>>>>>>")
	return nil, ErrNotSupport
}

// 默认编解码,只做长度和数据的解析
type DefaultCodec struct {
	RingBufferCodec
}

func NewDefaultCodec(packetHeaderSize PacketHeaderSize) *RingBufferCodec {
	if packetHeaderSize == nil {
		packetHeaderSize = func() uint32 {
			logger.Debug("  ------ring buff-------  PacketHeaderSize------")
			return uint32(DefaultPacketHeaderSize)
		}
	}
	codec := &RingBufferCodec{}
	codec.packetHeaderSize = packetHeaderSize
	return codec
}
