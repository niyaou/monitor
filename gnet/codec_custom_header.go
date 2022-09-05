package gnet

// proto.Message编解码
type CustomHeaderCodec struct {
	RingBufferCodec
	headerSize uint32
}

func NewCustomHeaderCodec(size uint32) *CustomHeaderCodec {
	ringBufferCodec := NewDefaultCodec(func() uint32 { return size })
	codec := &CustomHeaderCodec{
		RingBufferCodec: *ringBufferCodec,
	}
	codec.headerSize = size
	// codec.RingBufferCodec.PacketHeaderSize = codec.PacketHeaderSize
	return codec
}

func (_self *CustomHeaderCodec) PacketHeaderSize() uint32 {
	logger.Info("----custom--PacketHeaderSize-------%v", _self.headerSize)
	return uint32(_self.headerSize)
}
func (_self *CustomHeaderCodec) CreatePacketHeader(connection Connection, packet Packet, packetData []byte) PacketHeader {
	return NewDefaultPacketHeader(7, 0)
}
