package cluster

import (
	"go-redis/interface/resp"
	"go-redis/resp/reply"
)

// Rename 是通过key来进行hash计算的
// 需要判断是否落到同一个node中了
// rename k1 k2
func Rename(cluster *ClusterDatabase, c resp.Connection, cmdArgs [][]byte) resp.Reply {
	if len(cmdArgs) != 3 {
		return reply.MakeErrReply("Err wrong number args")
	}
	src := string(cmdArgs[1])
	dest := string(cmdArgs[2])
	// 源地址ip
	srcIP := cluster.peerPicker.PickNode(src)
	// 目标地址ip
	destIP := cluster.peerPicker.PickNode(dest)
	if srcIP != destIP {
		// TODO 删除掉旧的sec，在新的ip创建dest
		return reply.MakeErrReply("Err rename must within on IP")
	}
	return cluster.relay(srcIP, c, cmdArgs)
}
