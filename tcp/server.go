package tcp

import (
	"context"
	"go-redis/interface/tcp"
	"go-redis/lib/logger"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type Config struct {
	Address string
}

func ListenAndServeWithSignal(cfg *Config,
	handler tcp.Handler) error {
	// ctl connect
	closeChan := make(chan struct{})

	sigChan := make(chan os.Signal)
	// go runtime will detect this ... signal
	signal.Notify(sigChan, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)
	go func() {
		// 监听系统信号
		sig := <-sigChan
		switch sig {
		case syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL:
			closeChan <- struct{}{}
		}
	}()
	// listen address
	listener, err := net.Listen("tcp", cfg.Address)
	if err != nil {
		return err
	}
	logger.Info("start listen " + cfg.Address + " ...")
	ListenAndServe(listener, handler, closeChan)
	return nil
}

func ListenAndServe(
	listener net.Listener,
	handler tcp.Handler,
	closeChan <-chan struct{}) {
	// when kill process
	go func() {
		<-closeChan
		logger.Info("shutting down")
		_ = listener.Close()
		_ = handler.Close()

	}()

	defer func() {
		_ = listener.Close()
		_ = handler.Close()
	}()
	// 根据background create context
	ctx := context.Background()
	// 每连接一个godis客户端，都会让wg+1,当死循环出现异常的时候
	// 等待所有服务的客户做完工作后退出
	var waitDone sync.WaitGroup
	for true {
		conn, err := listener.Accept()
		if err != nil {
			// 接收新连接出现了问题
			logger.Error("Connection error ", err)
			break
		}
		logger.Info("accept new like")
		// 新增等待
		waitDone.Add(1)
		go func() {
			defer func() {
				waitDone.Done()
			}()
			handler.Handle(ctx, conn)
		}()
	}
	// 等待协程结束
	waitDone.Wait()
}
