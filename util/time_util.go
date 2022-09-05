package util

import (
	"fmt"
	"time"
)

// 获取当前时间戳(秒)
func GetCurrentTimeStamp() uint32 {
	return uint32(time.Now().UnixNano() / int64(time.Second))
}

// 获取当前毫秒数(毫秒,0.001秒)
func GetCurrentMS() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func GetCurrentTimeString() string {
	now := time.Now()
	return fmt.Sprintf("%02d-%02d-%02d-%02d-%02d-%02d", now.Year(), int(now.Month()),
		now.Day(), now.Hour(), now.Minute(), now.Second())
}
