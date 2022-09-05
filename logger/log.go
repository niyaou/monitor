package logger

import (
	"runtime"

	"desay.com/radar-monitor/gnet"
)

var logger = gnet.NewStdLogger(3)

func Debug(format string, args ...interface{}) {
	logger.Debug(format, args...)
}

func Info(format string, args ...interface{}) {
	logger.Info(format, args...)
}

func Warn(format string, args ...interface{}) {
	logger.Warn(format, args...)
}

func Error(format string, args ...interface{}) {
	logger.Error(format, args...)
}

func LogStack() {
	buf := make([]byte, 1<<12)
	Error(string(buf[:runtime.Stack(buf, false)]))
}
