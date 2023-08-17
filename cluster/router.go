package cluster

import "go-redis/interface/resp"

// 路由表
func makeRouter() map[string]CmdFunc {
	routerMap := make(map[string]CmdFunc)
	routerMap["exists"] = defaultFunc
	routerMap["type"] = defaultFunc
	routerMap["set"] = defaultFunc
	routerMap["setnx"] = defaultFunc
	routerMap["get"] = defaultFunc
	routerMap["getset"] = defaultFunc
	routerMap["ping"] = Ping
	routerMap["rename"] = Rename
	routerMap["renamenx"] = Rename
	routerMap["FlushDB"] = FlushDB
	routerMap["Del"] = Del
	routerMap["select"] = ExecSelect

	return routerMap
}

// 路由的默认方法，基本上都是在转发 get key : set k1 v1
func defaultFunc(cluster *ClusterDatabase, c resp.Connection, cmdArgs [][]byte) resp.Reply {
	key := string(cmdArgs[1])
	// 通过一致性hash 计算出来要转发的hash值
	peer := cluster.peerPicker.PickNode(key)
	// 直接进行转发
	return cluster.relay(peer, c, cmdArgs)
}
