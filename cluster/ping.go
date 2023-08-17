package cluster

import "go-redis/interface/resp"

// Ping 是本地执行模式
func Ping(cluster *ClusterDatabase, c resp.Connection, cmdArgs [][]byte) resp.Reply {
	return cluster.db.Exec(c, cmdArgs)
}
