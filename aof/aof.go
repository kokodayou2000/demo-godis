package aof

import (
	"go-redis/config"
	"go-redis/interface/database"
	"go-redis/lib/logger"
	"go-redis/lib/utils"
	"go-redis/resp/connection"
	"go-redis/resp/parser"
	"go-redis/resp/reply"
	"io"
	"os"
	"strconv"
)

type CmdLine [][]byte
type payload struct {
	cmdLine CmdLine
	dbIndex int
}

type AofHandler struct {
	database    database.Database
	aofChan     chan *payload
	aofFile     *os.File
	aofFilename string
	currentDB   int
}

const aofBufferSize = 1<<16 - 1

// NewAofHandler
func NewAofHandler(database database.Database) (*AofHandler, error) {

	handler := &AofHandler{}
	handler.aofFilename = config.Properties.AppendFilename
	handler.database = database
	handler.LoadAof()
	// perm 创建文件之后用什么权限打开 110 读写权限
	aofFile, err := os.OpenFile(handler.aofFilename, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}
	handler.aofFile = aofFile
	handler.aofChan = make(chan *payload, aofBufferSize)
	go func() {
		// 创建协程，执行落盘
		handler.handlerAof()
	}()
	return handler, err
}

// AddAof payload(set k v)
func (handler *AofHandler) AddAof(dbIndex int, cmdLine CmdLine) {
	if config.Properties.AppendOnly && handler.aofChan != nil {
		handler.aofChan <- &payload{
			cmdLine: cmdLine,
			dbIndex: dbIndex,
		}
	}

}

// handlerAof payload(set k v) <- aofCha(落盘)
// 需要注意的是，写的时候，需要标明在那个 storage 写的 比如 select 0
// *2 $5 select $1 3
// 但是每一个operator 上面都增加一个 select x 会很占空间
func (handler *AofHandler) handlerAof() {
	handler.currentDB = 0
	for p := range handler.aofChan {
		// currentDB记录的上次db号
		if p.dbIndex != handler.currentDB {
			// 添加select 0
			// ToCmdLine 将字符串转换成 [][]byte
			// toBytes 转换成 Redis 规范中的字符串 *2  $5 select ...
			data := reply.MakeMultiBulkReply(utils.ToCmdLine("select", strconv.Itoa(p.dbIndex))).ToBytes()
			_, err := handler.aofFile.Write(data)
			if err != nil {
				logger.Error(err)
				// 业务不能停止
				continue
			}
			// 记录是否切换db
			handler.currentDB = p.dbIndex
		}
		data := reply.MakeMultiBulkReply(p.cmdLine).ToBytes()
		_, err := handler.aofFile.Write(data)
		if err != nil {
			logger.Error(err)
			continue
		}
	}
}

// LoadAof 读盘
func (handler *AofHandler) LoadAof() {
	file, err := os.Open(handler.aofFilename)
	if err != nil {
		logger.Error(err)
		return
	}
	defer func() {
		err := file.Close()
		if err != nil {
			return
		}
	}()
	// 为了记录db index
	fakeConn := &connection.Connection{}
	// 将文件当成输入流，通过parse来解析
	ch := parser.ParseStream(file)
	for p := range ch {
		if p.Err != nil {
			// 如果读到文件结束符
			if p.Err == io.EOF {
				break
			}
			// 读取到的这一行指令有问题，打印异常，并继续
			logger.Error("aof文件读取异常" + p.Err.Error())
			continue
		}
		if p.Data == nil {
			// 读取到的这一行指令有问题，打印异常，并继续
			logger.Error("读取空指令" + p.Err.Error())
			continue
		}
		r, ok := p.Data.(*reply.MultiBulkReply)
		if !ok {
			logger.Error("读取的命令不是 multiBulkReply格式的")
		}
		rep := handler.database.Exec(fakeConn, r.Args)
		// 判断是否是错误的响应
		if reply.IsErrorReply(rep) {
			logger.Error(rep)
		}
	}
}
