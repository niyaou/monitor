package gnet

import (
	"context"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"google.golang.org/protobuf/proto"
)

// 连接接口定义
type Connection interface {
	// 连接唯一id
	GetConnectionId() uint32

	// 是否是发起连接的一方
	IsConnector() bool

	// 发包(protobuf)
	// NOTE:调用Send(command,message)之后,不要再对message进行读写!
	Send(command PacketCommand, message proto.Message) bool

	// 发包
	// NOTE:调用SendPacket(packet)之后,不要再对packet进行读写!
	SendPacket(packet Packet) bool

	// 超时发包,超时未发送则丢弃,适用于某些允许丢弃的数据包
	TrySendPacket(packet Packet, timeout time.Duration) bool

	// 是否连接成功
	IsConnected() bool

	// 获取编解码接口
	GetCodec() Codec

	// 设置编解码接口
	SetCodec(codec Codec)

	// LocalAddr returns the local network address.
	LocalAddr() net.Addr

	// RemoteAddr returns the remote network address.
	RemoteAddr() string

	SetAddr(addr string)

	// 关闭连接
	Close()

	// 获取关联数据
	GetTag() interface{}

	// 设置关联数据
	SetTag(tag interface{})

	Connect(address string) bool

	Start(ctx context.Context, netMgrWg *sync.WaitGroup, onClose func(connection Connection))
}

// 连接设置
type ConnectionConfig struct {
	// 发包缓存chan大小(缓存数据包chan容量)
	SendPacketCacheCap uint32
	// 发包Buffer大小(byte)
	// 不能小于PacketHeaderSize
	SendBufferSize uint32
	// 收包Buffer大小(byte)
	// 不能小于PacketHeaderSize
	RecvBufferSize uint32
	// 最大包体大小设置(byte),不包含PacketHeader
	// 允许该值大于SendBufferSize和RecvBufferSize
	MaxPacketSize uint32
	// 收包超时设置(秒)
	RecvTimeout uint32
	// 心跳包发送间隔(秒),对connector有效
	HeartBeatInterval uint32
	// 发包超时设置(秒)
	// net.Conn.SetWriteDeadline
	WriteTimeout uint32
	// TODO:其他流量控制设置
}

// 连接
type baseConnection struct {
	// 连接唯一id
	connectionId uint32
	// 连接ip端口
	addr string
	// 连接设置
	config *ConnectionConfig
	// 是否是连接方
	isConnector bool
	// 是否连接成功
	isConnected bool
	// 接口
	handler ConnectionHandler
	// 编解码接口
	codec Codec
	// 关联数据
	tag interface{}
}

// 连接唯一id
func (_self *baseConnection) GetConnectionId() uint32 {
	return _self.connectionId
}

func (_self *baseConnection) IsConnector() bool {
	return _self.isConnector
}

// 是否连接成功
func (_self *baseConnection) IsConnected() bool {
	return _self.isConnected
}

// 获取编解码接口
func (_self *baseConnection) GetCodec() Codec {
	return _self.codec
}

// 设置编解码接口
func (_self *baseConnection) SetCodec(codec Codec) {
	_self.codec = codec
}

// 获取关联数据
func (_self *baseConnection) GetTag() interface{} {
	return _self.tag
}

// 设置关联数据
func (_self *baseConnection) SetTag(tag interface{}) {
	_self.tag = tag
}

func (_self *baseConnection) SetAddr(addr string) {
	_self.addr = addr
}

var (
	connectionIdCounter uint32 = 0
)

func newConnectionId() uint32 {
	return atomic.AddUint32(&connectionIdCounter, 1)
}

type ConnectionCreator func(config *ConnectionConfig, address string, codec Codec, handler ConnectionHandler, isFixedHead bool) Connection

type AcceptConnectionCreator func(conn net.Conn, config *ConnectionConfig, codec Codec, handler ConnectionHandler) Connection
