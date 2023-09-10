package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/usiot/gbserver/handler"
	"github.com/usiot/gbserver/internal/config"
	"github.com/usiot/gbserver/internal/sip"
)

func main() {
	config.Init("goserve.json")
	handler.Init()
	sip.InitSip(&config.Conf.Sip)

	// 等待中断信号来优雅地关闭服务器，为关闭服务器操作设置一个5秒的超时
	quit := make(chan os.Signal, 1) // 创建一个接收信号的通道
	// kill 默认会发送 syscall.SIGTERM 信号
	// kill -2 发送 syscall.SIGINT 信号，我们常用的Ctrl+C就是触发系统SIGINT信号
	// kill -9 发送 syscall.SIGKILL 信号，但是不能被捕获，所以不需要添加它
	// signal.Notify把收到的 syscall.SIGINT或syscall.SIGTERM 信号转发给quit
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // 此处不会阻塞
	<-quit                                               // 阻塞在此，当接收到上述两种信号时才会往下执行
	log.Println("Shutdown Server ...")
	log.Println("Server exiting")
}
