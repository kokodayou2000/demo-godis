package connection

import (
	"go-redis/lib/logger"
	"go-redis/lib/sync/wait"
	"net"
	"sync"
	"time"
)

const (
	// flagSlave means this a connection with slave
	flagSlave = uint64(1 << iota)
	// flagSlave means this a connection with master
	flagMaster
	// flagMulti means this connection is within a transaction
	flagMulti
)

// Connection represents a connection with a redis-cli
type Connection struct {
	conn net.Conn

	// wait until finish sending data, used for graceful shutdown
	sendingData wait.Wait

	// lock while server sending response
	mu    sync.Mutex
	flags uint64

	// subscribing channels
	subs map[string]bool

	// password may be changed by CONFIG command during runtime, so store the password
	password string

	// queued commands for `multi`
	queue    [][][]byte
	watching map[string]uint32
	txErrors []error

	// selected db
	selectedDB int
}

var connPool = sync.Pool{
	New: func() interface{} {
		return &Connection{}
	},
}

// Close disconnect with the client
func (c *Connection) Close() error {
	c.sendingData.WaitWithTimeout(10 * time.Second)
	_ = c.conn.Close()
	c.subs = nil
	c.password = ""
	c.queue = nil
	c.watching = nil
	c.txErrors = nil
	c.selectedDB = 0
	connPool.Put(c)
	return nil
}

// NewConn creates Connection instance
func NewConn(conn net.Conn) *Connection {
	c, ok := connPool.Get().(*Connection)
	if !ok {
		logger.Error("connection pool make wrong type")
		return &Connection{
			conn: conn,
		}
	}
	c.conn = conn
	return c
}

// Write sends response to client over tcp connection
func (c *Connection) Write(b []byte) (int, error) {
	if len(b) == 0 {
		return 0, nil
	}
	c.sendingData.Add(1)
	defer func() {
		c.sendingData.Done()
	}()

	return c.conn.Write(b)
}

// GetDBIndex returns selected db
func (c *Connection) GetDBIndex() int {
	return c.selectedDB
}

func (c *Connection) SelectDBIndex(dbNum int) {
	c.selectedDB = dbNum
}
