package gnet

//测试服务之间利用通道传递消息，后续废弃

import (
	"sync"
)

// 网络管理类,提供对外接口
type MsgChan struct {
	_ch chan byte

	// 初始化一次
	initOnce sync.Once
}

var (
	// singleton
	msgChan = &MsgChan{}
)

// 单例模式,在调用的时候才会执行初始化一次
func GetMsgChan() *MsgChan {
	msgChan.initOnce.Do(func() {
		msgChan.init()
	})
	return msgChan
}

// 初始化
func (_self *MsgChan) init() {
	_self._ch = make(chan byte, 4)
}

// 初始化
func (_self *MsgChan) GetChan(ch *chan byte) *chan byte {
	return &_self._ch
}
