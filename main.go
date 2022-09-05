package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"

	"desay.com/radar-monitor/communication"
	"desay.com/radar-monitor/gnet"
	"desay.com/radar-monitor/internal"
	"desay.com/radar-monitor/logger"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			panic(err)
		}
	}()

	isDaemon := false

	// 配置文件名格式:
	flag.BoolVar(&isDaemon, "d", false, "daemon mode")
	flag.Parse()

	if isDaemon {
		daemon()
		return
	}

	runtime.GOMAXPROCS(runtime.NumCPU())
	gnet.SetLogLevel(gnet.DebugLevel)
	rand.Seed(time.Now().UnixNano())

	broker := internal.NewClient()
	broker.SetConditions(100)
	ctx, cancel := context.WithCancel(context.Background())
	serverList := []internal.Server{}

	// 根据命令行参数 创建不同的服务器实例
	// for _, serverType := range []string{"radar", "rotary", "rotary_module", "monitor"} {
	for _, serverType := range []string{"radar", "rotary", "rotary_module", "monitor"} {
		server := createServer(serverType)
		server.SetBroker(broker)
		if !server.Init(ctx, fmt.Sprintf("%s.json", serverType)) {
			panic("server init error")
		}
		// server.Run(ctx)
		// logger.Info("Run server %v", serverType)
		serverList = append(serverList, server)

	}
	logger.Info("serverList %v", serverList)

	// 监听系统的kill信号
	signalKillNotify := make(chan os.Signal, 1)
	signal.Notify(signalKillNotify, os.Interrupt, os.Kill, syscall.SIGTERM)
	if runtime.GOOS == "windows" {
		go func() {
			consoleReader := bufio.NewReader(os.Stdin)
			for {
				lineBytes, _, _ := consoleReader.ReadLine()
				line := strings.ToLower(string(lineBytes))
				logger.Debug("line:%v", line)
				if line == "close" || line == "exit" {
					logger.Debug("kill by console input")
					// 在windows系统模拟一个kill信号,以方便测试服务器退出流程
					signalKillNotify <- os.Kill
				}
			}
		}()
	}

	// 阻塞等待系统关闭信号
	logger.Info("wait for kill signal")
	select {
	case <-signalKillNotify:
		logger.Info("signalKillNotify, cancel ctx")
		// 通知所有协程关闭,所有监听<-ctx.Done()的地方会收到通知
		cancel()
		break
	}
	// 清理
	for _, server := range serverList {
		server.Exit()
	}
}

func daemon() {
	args := os.Args[1:]
	for i := 0; i < len(args); i++ {
		if args[i] == "-d=true" {
			args[i] = "-d=false"
			break
		}
	}
	cmd := exec.Command(os.Args[0], args...)
	cmd.Start()
	fmt.Println("[PID]", cmd.Process.Pid)
	os.Exit(0)
}

// 从配置文件名解析出服务器类型
func getServerTypeFromConfigFile(configFile string) string {
	baseFileName := filepath.Base(configFile)
	idx := strings.Index(baseFileName, "_")
	return baseFileName[0:idx]
}

// 创建相应类型的服务器
func createServer(serverType string) internal.Server {
	switch serverType {
	case "login":
		return nil
	case "game":
		return nil
	case "radar":
		return new(communication.CommunicateServer)
	case "monitor":
		return new(communication.MonitorViewServer)
	case "rotary":
		return new(communication.RotaryServer)
	case "rotary_module":
		return new(communication.RotaryModuleServer)
	}
	panic("err server type")
}
