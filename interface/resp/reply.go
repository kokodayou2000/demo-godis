package resp

// Reply 对客户端的回复
type Reply interface {
	ToBytes() []byte // tcp 通过字节传输数据
}
