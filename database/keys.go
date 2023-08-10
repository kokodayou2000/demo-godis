package database

import (
	"container/list"
	"go-redis/interface/resp"
	"go-redis/lib/wildcard"
	"go-redis/resp/reply"
)

func init() {
	RegisterCommand("del", execDel, -2)          // del k1 k2 k3
	RegisterCommand("exists", execExists, -2)    // exists k1 k2 k3
	RegisterCommand("flushDB", execFlushDB, -1)  // flushDB  使用 -1 不会后面的数进行纠错 用户输入 flushDB a b c 也不报错
	RegisterCommand("type", execType, 2)         // type k1
	RegisterCommand("rename", execRename, 3)     // rename k1 k2
	RegisterCommand("renameNX", execRenameNX, 3) // renameNX k1 k2
	RegisterCommand("keys", execKeys, 2)         // keys *

}

// DEL
func execDel(db *DB, args [][]byte) resp.Reply {
	keys := make([]string, len(args))
	for i, v := range args {
		keys[i] = string(v)
	}
	deleted := db.Removes(keys...)
	return reply.MakeIntReply(int64(deleted))
}

// EXISTS 检测key是否存在
func execExists(db *DB, args [][]byte) resp.Reply {
	result := int64(0)
	for _, arg := range args {
		key := string(arg)
		_, exists := db.GetEntity(key)
		if exists {
			result++
		}
	}
	return reply.MakeIntReply(result)
}

// FlushDB
func execFlushDB(db *DB, args [][]byte) resp.Reply {
	db.Flush()
	return reply.MakeOKReply()
}

// TYPE
func execType(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	entity, exists := db.GetEntity(key)
	if !exists {
		return reply.MakeStatusReply("none")
	}
	switch entity.Data.(type) {
	case []byte:
		return reply.MakeStatusReply("string")
	case list.List:
		return reply.MakeStatusReply("list")
	}
	return &reply.UnKnowErrReply{}
}

// Rename k1 k2
func execRename(db *DB, args [][]byte) resp.Reply {
	src := string(args[0])
	dest := string(args[1])
	entity, exist := db.GetEntity(src)
	if !exist {
		return reply.MakeErrReply("no such key")
	}
	db.PutEntity(dest, entity)
	db.Remove(src)
	return reply.MakeOKReply()
}

// RenameNX
// 返回值和rename不同，rename会返回ok，但是renameNX未进行操作返回0，进行操作返回1
func execRenameNX(db *DB, args [][]byte) resp.Reply {
	src := string(args[0])
	dest := string(args[1])
	_, exist := db.GetEntity(dest)
	if exist {
		// redis 中如果没有进行操作，就返回 0
		return reply.MakeIntReply(0)
	}

	entity, exist := db.GetEntity(src)
	if !exist {
		return reply.MakeErrReply("no such key")
	}
	db.PutEntity(dest, entity)
	db.Remove(src)
	return reply.MakeIntReply(1)
}

// KEYS * redis 通配符算法
func execKeys(db *DB, args [][]byte) resp.Reply {
	// 根据 参数获取对应的 regex
	pattern, _ := wildcard.CompilePattern(string(args[0]))
	result := make([][]byte, 0)
	db.data.ForEach(func(key string, val interface{}) bool {
		// 根据regex 匹配key，判断是否能匹配成功
		if pattern.IsMatch(key) {
			result = append(result, []byte(key))
		}
		return true
	})
	return reply.MakeMultiBulkReply(result)
}
