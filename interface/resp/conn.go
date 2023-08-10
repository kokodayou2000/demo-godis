package resp

// Connection 代表redis协议层的连接
type Connection interface {
	Write([]byte) (int, error)
	GetDBIndex() int   // default 16 DB
	SelectDBIndex(int) // will change db connection
}
