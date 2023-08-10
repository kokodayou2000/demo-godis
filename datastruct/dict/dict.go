package dict

// Consumer 定义消费方法类型
type Consumer func(key string, val interface{}) bool

// Dict 字典类型接口
type Dict interface {
	// Get 通过key获取到val，如果存在返回val，如果不存在返回 false
	Get(key string) (val interface{}, exist bool)
	// Len 字典的数据量
	Len() int
	// Put  存放 key value
	Put(key string, val interface{}) (result int)
	// PutIfAbsent such as setNX
	PutIfAbsent(key string, val interface{}) (result int)
	// PutIfExists such as setEX
	PutIfExists(key string, val interface{}) (result int)
	// Remove 删除key
	Remove(key string) (result int)
	// ForEach 遍历整个godis key value
	// 把整个godis数据当成一个stream foreach
	ForEach(consumer Consumer)
	// Keys 获取key集合
	Keys() []string
	// RandomKeys 随机获取limit个key
	RandomKeys(limit int) []string
	// RandomDistinctKeys 随机返回不重复的key
	RandomDistinctKeys(limit int) []string
	// Clear flushDB
	Clear()
}
