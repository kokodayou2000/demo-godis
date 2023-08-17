package cluster

import (
	"context"
	pool "github.com/jolestar/go-commons-pool/v2"
	"go-redis/config"
	database2 "go-redis/database"
	"go-redis/interface/database"
	"go-redis/interface/resp"
	"go-redis/lib/consistenthash"
	"go-redis/lib/logger"
	"go-redis/resp/reply"
	"strings"
)

type ClusterDatabase struct {
	self           string
	nodes          []string
	peerPicker     *consistenthash.NodeMap     //节点选择器
	peerConnection map[string]*pool.ObjectPool //使用map来维护多个连接池
	db             database.Database           // 集群层从下一层 standalone
}

func MakeClusterDatabase() *ClusterDatabase {
	cd := &ClusterDatabase{
		self:           config.Properties.Self,
		db:             database2.NewStandaloneDatabase(),
		peerPicker:     consistenthash.NewNodeMap(nil),
		peerConnection: make(map[string]*pool.ObjectPool),
	}
	// 解析配置文件中的节点，放到数组中
	// 节点的长度 capable =  len(config.Properties.Peers)+1
	nodes := make([]string, 0, len(config.Properties.Peers)+1)
	for _, peer := range config.Properties.Peers {
		// 放到数组中
		nodes = append(nodes, peer)
	}
	nodes = append(nodes, config.Properties.Self)
	// add
	cd.peerPicker.AddNode(nodes...)
	// 为每一个创建连接池
	ctx := context.Background()
	for _, peer := range config.Properties.Peers {
		// 为每一个peer创建一个连接池
		cd.peerConnection[peer] = pool.NewObjectPoolWithDefaultConfig(ctx, &connectionFactory{Peer: peer})
	}
	cd.nodes = nodes
	return cd
}

// CmdFunc 规定一个在集群中执行的方法
type CmdFunc func(cluster *ClusterDatabase, c resp.Connection, cmdArgs [][]byte) resp.Reply

var router = makeRouter()

// Exec 集群版本的执行 (result resp.Reply) 相当于在第一行创建了一个变量 result resp.Reply 语法糖
func (cd *ClusterDatabase) Exec(client resp.Connection, args [][]byte) (result resp.Reply) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error(err)
			result = &reply.UnknownErrReply{}
		}
	}()
	cmdName := strings.ToLower(string(args[0]))
	cmdFunc, ok := router[cmdName]
	if !ok {
		return reply.MakeErrReply("not supported cmd")
	}
	// 执行
	result = cmdFunc(cd, client, args)
	return
}

func (cd *ClusterDatabase) AfterClientClose(conn resp.Connection) {
	cd.db.AfterClientClose(conn)
}

func (cd *ClusterDatabase) Close() {
	cd.db.Close()
}
