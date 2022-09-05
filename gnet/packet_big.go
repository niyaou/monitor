package gnet

import (
	"encoding/binary"
	"unsafe"

	"google.golang.org/protobuf/proto"
)

const (
	// 大包的包头长度
	DefaultBigPacketHeaderSize = int(unsafe.Sizeof(BigPacketHeader{}))
	// 数据包长度限制(4G)
	MaxBigPacketDataSize = 0xFFFFFFFF
)

// 大包包头
// 大包:包长度可能超过16M
type BigPacketHeader struct {
	len     uint32
	command uint16
	flags   uint16
}

func NewBigPacketHeader(len uint32, command, flags uint16) *BigPacketHeader {
	return &BigPacketHeader{
		len:     len,
		command: command,
		flags:   flags,
	}
}

// 包体长度,不包含包头的长度
// [0,0xFFFFFFFF]
func (_self *BigPacketHeader) Len() uint32 {
	return _self.len
}

// 消息号
func (_self *BigPacketHeader) Command() uint16 {
	return _self.command
}

// 标记
func (_self *BigPacketHeader) Flags() uint16 {
	return _self.flags
}

// 从字节流读取数据,len(messageHeaderData)>=MessageHeaderSize
// 使用小端字节序
func (_self *BigPacketHeader) ReadFrom(messageHeaderData []byte) {
	_self.len = binary.LittleEndian.Uint32(messageHeaderData)
	_self.command = binary.LittleEndian.Uint16(messageHeaderData[4:])
	_self.flags = binary.LittleEndian.Uint16(messageHeaderData[6:])
}

// 写入字节流,使用小端字节序
func (_self *BigPacketHeader) WriteTo(messageHeaderData []byte) {
	binary.LittleEndian.PutUint32(messageHeaderData, _self.len)
	binary.LittleEndian.PutUint16(messageHeaderData[4:], _self.command)
	binary.LittleEndian.PutUint16(messageHeaderData[6:], _self.flags)
}

// 包含一个消息号和[]byte的数据包
type BigDataPacket struct {
	command uint16
	data    []byte
}

func NewBigDataPacket(command uint16, data []byte) *BigDataPacket {
	return &BigDataPacket{
		command: command,
		data:    data,
	}
}

func (_self *BigDataPacket) Command() PacketCommand {
	return PacketCommand(_self.command)
}

func (_self *BigDataPacket) Message() proto.Message {
	return nil
}

func (_self *BigDataPacket) GetStreamData() []byte {
	return _self.data
}

func (_self *BigDataPacket) GetStringData() string {
	return ""
}

// deep copy
func (_self *BigDataPacket) Clone() Packet {
	newPacket := &BigDataPacket{data: make([]byte, len(_self.data))}
	newPacket.command = _self.command
	copy(newPacket.data, _self.data)
	return newPacket
}
