package cluster

import (
	"go-redis/interface/resp"
	"go-redis/resp/reply"
)

// Del k1 k2 k3 k4 可能出现在不同的节点
func Del(cluster *ClusterDatabase, c resp.Connection, cmdArgs [][]byte) resp.Reply {
	replies := cluster.broadcast(c, cmdArgs)
	var errReply reply.ErrorReply
	var deleted int64 = 0
	for _, r := range replies {
		if reply.IsErrorReply(r) {
			// 只要有一个node没有执行flushDB成功
			errReply = r.(reply.ErrorReply)
			break
		}
		intReply, ok := r.(*reply.IntReply)
		if !ok {
			errReply = reply.MakeErrReply("Del error")
		}
		deleted += intReply.Code
	}
	if errReply == nil {
		// 返回删除了多少个
		return reply.MakeIntReply(deleted)
	}
	return reply.MakeErrReply("error:" + errReply.Error())
}
