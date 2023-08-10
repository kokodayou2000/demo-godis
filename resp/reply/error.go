package reply

type UnKnowErrReply struct {
}

var unKnowErrorBytes = []byte("-Err unknown\r\n")

func (u UnKnowErrReply) Error() string {
	return "Err unknown"
}

func (u UnKnowErrReply) ToBytes() []byte {
	return unKnowErrorBytes
}

type ArgNumErrReply struct {
	Cmd string
}

func (arg *ArgNumErrReply) Error() string {
	return "-ERR wrong number of arguments for '" + arg.Cmd + "' command\r\n"
}

func (arg *ArgNumErrReply) ToBytes() []byte {
	return []byte("-ERR wrong number of arguments for '" + arg.Cmd + "' command\r\n")
}

func MakeArgNumErrReply(cmd string) *ArgNumErrReply {
	return &ArgNumErrReply{Cmd: cmd}
}

type SyntaxErrReply struct {
}

var syntaxErrBytes = []byte("-Err syntax error\r\n")
var theSyntaxErrReply = &SyntaxErrReply{}

func MakeSyntaxErrReply() *SyntaxErrReply {
	return theSyntaxErrReply
}
func (s *SyntaxErrReply) ToBytes() []byte {
	return syntaxErrBytes
}
func (r *SyntaxErrReply) Error() string {
	return "Err syntax error"
}

// WrongTypeErrReply 数据类型错误
type WrongTypeErrReply struct {
}

var wrongTypeErrBytes = []byte("-WRONGTYPE Operation against a key holding the value\r\n")

func (s *WrongTypeErrReply) ToBytes() []byte {
	return wrongTypeErrBytes
}
func (r *WrongTypeErrReply) Error() string {
	return "WRONGTYPE Operation against a key holding the value"
}

// ProtocolErrReply 不符合RESP规范
type ProtocolErrReply struct {
	Msg string
}

func (r *ProtocolErrReply) ToBytes() []byte {
	return []byte("-ERR Protocol error: '" + r.Msg + "'\r\n")
}

func (r *ProtocolErrReply) Error() string {
	return "ERR Protocol error: " + r.Msg
}
