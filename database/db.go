package database

import (
	"go-redis/datastruct/dict"
	"go-redis/interface/database"
	"go-redis/interface/resp"
	"go-redis/resp/reply"
	"strings"
)

// DB 对应redis中16个槽
type DB struct {
	// 索引
	index int
	// 数据
	data dict.Dict
}

type CmdLine = [][]byte

type ExecFunc func(db *DB, args [][]byte) resp.Reply

// MakeDB 创建
func MakeDB() *DB {
	db := &DB{
		data: dict.MakeSyncDict(),
	}
	return db
}

// Exec 执行
// 1.连接信息  2.需要执行的命令行
func (db *DB) Exec(c resp.Connection, cmdLine CmdLine) resp.Reply {
	// example, set setNX get
	cmdName := strings.ToLower(string(cmdLine[0]))

	cmd, ok := cmdTable[cmdName]

	if !ok {
		// 不存在该命令
		return reply.MakeErrReply("ERR unknown command" + cmdName)
	}
	if !validateArity(cmd.arity, cmdLine) {
		// 参数个数校验
		return reply.MakeArgNumErrReply(cmdName)
	}
	// get execute set func
	// 获取到执行set方法的方法
	fun := cmd.executor
	// 只需要传入 set key value 中的key value即可
	return fun(db, cmdLine[1:])
}

func validateArity(arity int, cmd CmdLine) bool {
	return true
}

// GetEntity
func (db *DB) GetEntity(key string) (*database.DataEntity, bool) {
	raw, ok := db.data.Get(key)
	if !ok {
		return nil, false
	}
	// 对数据进行转换了
	entity, _ := raw.(*database.DataEntity)
	return entity, true
}

func (db *DB) PutEntity(key string, entity *database.DataEntity) int {
	return db.data.Put(key, entity)
}

// PutIfAbsent 只有不存在的时候，才添加
func (db *DB) PutIfAbsent(key string, entity *database.DataEntity) int {
	return db.data.PutIfAbsent(key, entity)
}

func (db *DB) Remove(key string) {
	db.data.Remove(key)
}

// Removes the given keys from db
func (db *DB) Removes(keys ...string) (deleted int) {
	deleted = 0
	for _, key := range keys {
		_, exist := db.data.Get(key)
		if exist {
			db.Remove(key)
			deleted++
		}
	}
	return deleted
}

func (db *DB) Flush() {
	db.data.Clear()
}
