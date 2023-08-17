package cluster

import "go-redis/interface/resp"

// ExecSelect 是本地执行模式
func ExecSelect(cluster *ClusterDatabase, c resp.Connection, cmdArgs [][]byte) resp.Reply {
	// 本地执行
	return cluster.db.Exec(c, cmdArgs)
}
