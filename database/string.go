package database

import (
	"go-redis/interface/database"
	"go-redis/interface/resp"
	"go-redis/resp/reply"
)

func init() {
	RegisterCommand("get", execGet, 2)
	RegisterCommand("set", execSet, 3)
	RegisterCommand("setNX", execSetNX, 3)
	RegisterCommand("getSet", execGetSet, 3)
	RegisterCommand("strlen", execStrLen, 2)
}

// Get
func execGet(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	entity, exist := db.GetEntity(key)
	if !exist {
		return reply.MakeNullBulkReply()
	}
	bytes, ok := entity.Data.([]byte)
	if ok {
		return reply.MakeBulkReply(bytes)
	}
	return reply.MakeErrReply("It doesn't string type")
}

// Set k1 v1
func execSet(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	value := string(args[1])
	// 为什么要使用 DataEntity ?
	entity := &database.DataEntity{Data: value}
	db.PutEntity(key, entity)
	return reply.MakeOKReply()
}

// SetNX k1 v1 if k1 exist ,return 0
func execSetNX(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	value := string(args[1])
	// 为什么要使用 DataEntity ?
	entity := &database.DataEntity{Data: value}
	result := db.PutIfAbsent(key, entity)
	return reply.MakeIntReply(int64(result))
}

// GetSet k1 v1 返回k1 原来的值，然后k1重新赋值
func execGetSet(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	value := string(args[1])
	oldEntity, exists := db.GetEntity(key)
	// 为什么要使用 DataEntity ?
	entity := &database.DataEntity{Data: value}
	_ = db.PutEntity(key, entity)
	if !exists {
		return reply.MakeNullBulkReply()
	}
	return reply.MakeBulkReply(oldEntity.Data.([]byte))
}

// strlen
func execStrLen(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	entity, exists := db.GetEntity(key)
	if !exists {
		return reply.MakeNullBulkReply()
	}
	bytes := entity.Data.([]byte)
	return reply.MakeIntReply(int64(len(bytes)))
}
