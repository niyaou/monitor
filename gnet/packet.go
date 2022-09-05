package gnet

import (
	"encoding/binary"
	"unsafe"

	"google.golang.org/protobuf/proto"
)

const (
	// 包头长度
	DefaultPacketHeaderSize = int(unsafe.Sizeof(DefaultPacketHeader{}))
	// 数据包长度限制
	MaxPacketDataSize = 0x00FFFFFF
)

// 消息号
type PacketCommand uint16

// 包头接口
type PacketHeader interface {
	Len() uint32
	ReadFrom(messageHeaderData []byte)
	WriteTo(messageHeaderData []byte)
}

// 默认包头,支持小于16M的数据包
type DefaultPacketHeader struct {
	// (flags << 24) | len
	// flags [0,255)
	// len [0,16M)
	LenAndFlags uint32
}

func NewDefaultPacketHeader(len uint32, flags uint8) *DefaultPacketHeader {
	return &DefaultPacketHeader{
		LenAndFlags: uint32(flags)<<24 | len,
	}
}

// func NewPacketHeader(len uint32, flags uint8) *PacketHeader {
// 	return &PacketHeader{
// 		LenAndFlags: uint32(flags)<<24 | len,
// 	}
// }

// 包体长度,不包含包头的长度
// [0,0x00FFFFFF]
func (_self *DefaultPacketHeader) Len() uint32 {
	return _self.LenAndFlags & 0x00FFFFFF
}

// 标记 [0,0xFF]
func (_self *DefaultPacketHeader) Flags() uint32 {
	return _self.LenAndFlags >> 24
}

// 从字节流读取数据,len(messageHeaderData)>=MessageHeaderSize
// 使用小端字节序
func (_self *DefaultPacketHeader) ReadFrom(messageHeaderData []byte) {
	_self.LenAndFlags = binary.LittleEndian.Uint32(messageHeaderData)
}

// 写入字节流,使用小端字节序
func (_self *DefaultPacketHeader) WriteTo(messageHeaderData []byte) {
	binary.LittleEndian.PutUint32(messageHeaderData, _self.LenAndFlags)
}

// 数据包接口
type Packet interface {
	// 消息号
	// 没有把消息号放在PacketHeader里,因为对TCP网络层来说,只需要知道每个数据包的分割长度就可以了,
	// 至于数据包具体的格式,不该是网络层关心的事情
	// 消息号也不是必须放在这里的,但是游戏项目一般都是用消息号,为了减少封装层次,就放这里了
	Command() PacketCommand

	// 默认使用protobuf
	Message() proto.Message

	// 预留一个二进制数据的接口,支持外部直接传入序列号的字节流数据
	GetStreamData() []byte

	//字符串数据流
	GetStringData() string

	// deep copy
	Clone() Packet
}

// proto数据包
type ProtoPacket struct {
	command       PacketCommand
	message       proto.Message
	data          []byte
	stringPayload string
}

func NewProtoPacket(command PacketCommand, message proto.Message) *ProtoPacket {
	return &ProtoPacket{
		command: command,
		message: message,
	}
}

func NewProtoPacketWithData(command PacketCommand, data []byte) *ProtoPacket {
	return &ProtoPacket{
		command: command,
		data:    data,
	}
}

func NewStringPacketWithData(command PacketCommand, payload string) *ProtoPacket {
	return &ProtoPacket{
		command:       command,
		stringPayload: payload,
	}
}

func (_self *ProtoPacket) Command() PacketCommand {
	return _self.command
}

func (_self *ProtoPacket) Message() proto.Message {
	return _self.message
}

// 某些特殊需求会直接使用序列化好的数据
func (_self *ProtoPacket) GetStreamData() []byte {
	return _self.data
}

func (_self *ProtoPacket) GetStringData() string {
	return _self.stringPayload
}

// deep copy
func (_self *ProtoPacket) Clone() Packet {
	return &ProtoPacket{
		command: _self.command,
		message: proto.Clone(_self.message),
	}
}

// 只包含一个[]byte的数据包
type DataPacket struct {
	data          []byte
	stringPayload string
}

func NewDataPacket(data []byte, payload string) *DataPacket {

	_pack := &DataPacket{data: data}
	if payload != "" {
		_pack.stringPayload = payload
	}
	return _pack
}

func (_self *DataPacket) Command() PacketCommand {
	return 0
}

func (_self *DataPacket) Message() proto.Message {
	return nil
}

func (_self *DataPacket) GetStreamData() []byte {
	return _self.data
}

func (_self *DataPacket) GetStringData() string {
	return _self.stringPayload
}

// deep copy
func (_self *DataPacket) Clone() Packet {
	newPacket := &DataPacket{data: make([]byte, len(_self.data))}
	copy(newPacket.data, _self.data)
	return newPacket
}

// type StringPacket struct {
// 	stringPayload string
// }

// func NewStringPacket(data string) *StringPacket {
// 	return &StringPacket{stringPayload: data}
// }

// func (_self *StringPacket) Command() PacketCommand {
// 	return 0
// }

// func (_self *StringPacket) Message() proto.Message {
// 	return nil
// }

// func (_self *StringPacket) GetStreamData() []byte {
// 	return nil
// }

// func (_self *StringPacket) GetStringData() string {
// 	return _self.stringPayload
// }

// // deep copy
// func (_self *StringPacket) Clone() Packet {
// 	newPacket := &StringPacket{stringPayload: _self.stringPayload}
// 	return newPacket
// }
