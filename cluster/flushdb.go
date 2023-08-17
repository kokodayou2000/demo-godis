package cluster

import (
	"go-redis/interface/resp"
	"go-redis/resp/reply"
)

func FlushDB(cluster *ClusterDatabase, c resp.Connection, cmdArgs [][]byte) resp.Reply {
	replies := cluster.broadcast(c, cmdArgs)
	var errReply reply.ErrorReply
	for _, r := range replies {
		if reply.IsErrorReply(r) {
			// 只要有一个node没有执行flushDB成功
			errReply = r.(reply.ErrorReply)
			break
		}
	}
	if errReply == nil {
		return reply.MakeOkReply()
	}
	return reply.MakeErrReply("error:" + errReply.Error())
}
