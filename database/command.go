package database

import "strings"

// 记录系统中所有指令和command的关系
var cmdTable = make(map[string]*command)

// command 命令
type command struct {
	// 执行的方法
	executor ExecFunc
	// 参数数量
	arity int
}

func RegisterCommand(name string, executor ExecFunc, arity int) {
	name = strings.ToLower(name)
	cmdTable[name] = &command{executor: executor, arity: arity}
}
