package handler

import (
	"context"
	"go-redis/database"
	dbface "go-redis/interface/database"
	"go-redis/lib/sync/atomic"
	"go-redis/resp/connection"
	"net"
	"sync"
)

var (
	unknownErrReplyBytes = []byte("-ERR unknown\r\n")
)

type RespHandler struct {
	activeConn sync.Map
	db         dbface.Database
	closing    atomic.Boolean
}

func MakeHandler() *RespHandler {
	var db dbface.Database
	db = database.NewEchoDatabase()
	return &RespHandler{db: db}
}

func (r *RespHandler) closeClient(client *connection.Connection) {
	_ = client.Close()
	r.db.AfterClientClose(client)
	r.activeConn.Delete(client)
}

func (r *RespHandler) Handler(ctx context.Context, conn net.Conn) {
	if r.closing.Get() {
		_ = conn.Close()
	}
}
