package cluster

import (
	"context"
	"errors"
	"go-redis/interface/resp"
	"go-redis/lib/utils"
	"go-redis/resp/client"
	"go-redis/resp/reply"
	"strconv"
)

// 通信

// 在连接池中拿到一个连接
func (cd *ClusterDatabase) getPeerClient(peer string) (*client.Client, error) {
	// 获取连接池
	objectPool, ok := cd.peerConnection[peer]
	if !ok {
		return nil, errors.New("connect not found")
	}
	// 获取连接池中的对象
	object, err := objectPool.BorrowObject(context.Background())
	if err != nil {
		return nil, err
	}
	c, ok := object.(*client.Client)
	if !ok {
		return nil, errors.New("pool type error")
	}
	return c, nil
}

// 把客户端还给连接池
func (cd *ClusterDatabase) returnPeerClient(peer string, client *client.Client) error {
	pool, ok := cd.peerConnection[peer]
	if !ok {
		return errors.New("连接池中不存在该 地址")
	}
	return pool.ReturnObject(context.Background(), client)
}

// 转发
// peer 目标节点，连接对象，参数
func (cd *ClusterDatabase) relay(peer string, conn resp.Connection, args [][]byte) resp.Reply {
	// 如果目标是自己
	if peer == cd.self {
		return cd.db.Exec(conn, args)
	}
	peerClient, err := cd.getPeerClient(peer)
	if err != nil {
		return reply.MakeErrReply(err.Error())
	}
	defer func() {
		_ = cd.returnPeerClient(peer, peerClient)
	}()
	// 先转发用户选择的db
	peerClient.Send(utils.ToCmdLine("select", strconv.Itoa(conn.GetDBIndex())))
	// send cmdLine
	return peerClient.Send(args)
}

// 需要广播的命令，比如 flushDB,将集群中的所有节点清库,会返回多个nodeName 以及对应的 reply
func (cd *ClusterDatabase) broadcast(conn resp.Connection, args [][]byte) map[string]resp.Reply {
	results := make(map[string]resp.Reply)
	for _, node := range cd.nodes {
		relay := cd.relay(node, conn, args)
		results[node] = relay
	}
	return results
}
