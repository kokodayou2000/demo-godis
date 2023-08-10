package database

import "go-redis/interface/resp"

type Database interface {
	Exec(client resp.Connection, cmdLine [][]byte) resp.Reply
	Close()
	AfterClientClose(c resp.Connection)
}

type DataEntity struct {
	Data interface{}
}
