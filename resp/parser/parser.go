package parser

import (
	"go-redis/interface/resp"
	"io"
)

// Payload 客户发送的数据
type Payload struct {
	Data resp.Reply
	Err  error
}

// readState 解析器的状态
type readState struct {
	// 是否多行
	readingMultiLine bool
	// 读取的指令应该有多少个参数
	expectedArgsCount int
	//消息类型
	msgType byte
	// 已经解析的参数
	args [][]byte
	// 整个字节的长度
	bulkLen int64
}

// 计算解析器是否完成
// 根据readState的 expectedArgsCount 来判定
func (s *readState) finished() bool {
	//期望的参数大于0，而且我们读取到了足够的参数，就表示解析成功
	return s.expectedArgsCount > 0 && len(s.args) == s.expectedArgsCount
}

// ParseStream 协议层对外的接口
// tcp 服务层调用将io流通过 ParseStream 来解析
// 并通过管道的方式将解析结果返回 (异步)
func ParseStream(reader io.Reader) <-chan *Payload {
	// 创建一个管道
	ch := make(chan *Payload)
	// 使用协程的方式来异步处理
	go parse0(reader, ch)

	return ch
}

// 具体执行解析的，从tcp成读取io流
func parse0(reader io.Reader, ch chan<- *Payload) {

}
