package reply

// PongReply  ping-pong
type PongReply struct {
}

var pongBytes = []byte("+PONG\r\n")

// ToBytes 将指定的 reply 返回
func (p *PongReply) ToBytes() []byte {
	return pongBytes
}

func MakePongReply() *PongReply {
	return &PongReply{}
}

// PongReply  ping-pong
type OKReply struct {
}

var okBytes = []byte("+OK\r\n")

// ToBytes 将指定的 reply 返回
func (p *OKReply) ToBytes() []byte {
	return okBytes
}

func MakeOKReply() *PongReply {
	return &PongReply{}
}

type NullBulkReply struct {
}

var nullBulkBytes = []byte("$-1\r\n")

func (n *NullBulkReply) ToBytes() []byte {
	return nullBulkBytes
}

func MakeNullBulkReply() *NullBulkReply {
	return &NullBulkReply{}
}

var emptyMultiBulkBytes = []byte("*0\r\n")

type EmptyMultiBulkReply struct {
}

func (e *EmptyMultiBulkReply) ToBytes() []byte {
	return emptyMultiBulkBytes
}

var noBytes = []byte("")

type NoReply struct {
}

func (n *NoReply) ToBytes() []byte {
	return noBytes
}
