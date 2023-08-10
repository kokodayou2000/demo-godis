package database

import (
	"go-redis/interface/resp"
	"go-redis/resp/reply"
)

func Ping(db *DB, args [][]byte) resp.Reply {
	//
	return reply.MakePongReply()
}

// store ping
func init() {
	RegisterCommand("ping", Ping, 1)
}
