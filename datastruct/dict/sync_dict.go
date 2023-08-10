package dict

import (
	"sync"
)

// SyncDict 并发安全的字典底层数据结构
type SyncDict struct {
	m sync.Map
}

// MakeSyncDict 创建
func MakeSyncDict() *SyncDict {
	return &SyncDict{}
}

func (dict *SyncDict) Get(key string) (val interface{}, exist bool) {
	value, ok := dict.m.Load(key)
	return value, ok
}

func (dict *SyncDict) Len() int {
	lenth := 0
	dict.m.Range(func(key, value any) bool {
		lenth++
		return true
	})
	return lenth
}

func (dict *SyncDict) Put(key string, val interface{}) (result int) {
	_, existed := dict.m.Load(key)
	dict.m.Store(key, val)
	if !existed {
		// insert
		return 1
	}
	// just update
	return 0

}

func (dict *SyncDict) PutIfAbsent(key string, val interface{}) (result int) {
	_, existed := dict.m.Load(key)
	if !existed {
		dict.m.Store(key, val)
		return 1
	}
	return 0
}

func (dict *SyncDict) PutIfExists(key string, val interface{}) (result int) {
	_, existed := dict.m.Load(key)
	if !existed {
		dict.m.Store(key, val)
		return 1
	}
	return 0
}

func (dict *SyncDict) Remove(key string) (result int) {
	_, existed := dict.m.Load(key)
	dict.m.Delete(key)
	if !existed {
		return 0
	}
	return 1
}

// ForEach 不做consumer结束判断，完整的遍历整个字典
func (dict *SyncDict) ForEach(consumer Consumer) {
	dict.m.Range(func(key, value any) bool {
		consumer(key.(string), value)
		return true
	})
}

// Keys 返回keys切片
func (dict *SyncDict) Keys() []string {
	keys := make([]string, dict.Len())
	len := 0
	dict.m.Range(func(key, value any) bool {
		keys[len] = key.(string)
		len++
		return true
	})
	return keys
}

func (dict *SyncDict) RandomKeys(limit int) []string {
	keys := make([]string, dict.Len())
	for i := 0; i < limit; i++ {
		dict.m.Range(func(key, value any) bool {
			keys[i] = key.(string)
			return false
		})
	}
	return keys
}

func (dict *SyncDict) RandomDistinctKeys(limit int) []string {
	keys := make([]string, dict.Len())
	i := 0
	// 利用本身的随机性
	dict.m.Range(func(key, value any) bool {
		if i <= limit {
			keys[i] = key.(string)
			i++
			return true
		}
		return false
	})
	return keys
}

func (dict *SyncDict) Clear() {
	// 新建一个dict
	*dict = *MakeSyncDict()
}
