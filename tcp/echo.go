package tcp

import (
	"bufio"
	"context"
	"go-redis/lib/logger"
	"go-redis/lib/sync/atomic"
	"go-redis/lib/sync/wait"
	"io"
	"net"
	"sync"
	"time"
)

// EchoClient 客户端实体
type EchoClient struct {
	Conn    net.Conn
	Waiting wait.Wait
}

// Close 关闭客户端的逻辑
func (e *EchoClient) Close() error {
	// 等待10s，如果还未能断开连接就关闭
	e.Waiting.WaitWithTimeout(10 * time.Second)
	_ = e.Conn.Close()
	return nil
}

// EchoHandler  业务 接收客户端连接并服务，以及关闭
type EchoHandler struct {
	activeConn sync.Map
	closing    atomic.Boolean // closing state
}

func NewHandler() *EchoHandler {
	// 使用 sync.Map和atomic.Boolean的默认初始值就好
	return &EchoHandler{}
}

func (handler *EchoHandler) Handle(ctx context.Context, conn net.Conn) {
	//我们的业务是否处于正在关闭的状态
	if handler.closing.Get() {
		_ = conn.Close()
	}
	// create client
	client := &EchoClient{
		Conn: conn,
	}
	// 存储到sync.Map中
	handler.activeConn.Store(client, struct{}{})
	// 创建连接对应的buffer
	reader := bufio.NewReader(conn)
	for true {
		// 结束的标志
		msg, err := reader.ReadString('\n')
		if err != nil {
			// 网络传输结束,客户端退出
			if err == io.EOF {
				logger.Info("Connecting Close...")
				handler.activeConn.Delete(client)
			} else {
				logger.Warn("Warn", err)
			}
			return
		}
		// 在wg中做一个记录，不要直接关闭我，等10s，或者等到Done之后
		client.Waiting.Add(1)
		b := []byte(msg)
		// 写回去
		_, _ = conn.Write(b)

		// 业务结束
		client.Waiting.Done()
	}
}

func (handler *EchoHandler) Close() error {
	logger.Info("handler shutting down")
	// 进入关闭状态
	handler.closing.Set(true)
	// close all connection
	handler.activeConn.Range(func(key, value any) bool {
		// 转换成client类型
		client := key.(*EchoClient)
		// Close
		_ = client.Conn.Close()
		return true
	})
	return nil
}
