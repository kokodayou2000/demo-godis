package consistenthash

import (
	"hash/crc32"
	"sort"
)

// HashFunc 参数是字节，返回值是落点
type HashFunc func(data []byte) uint32

type NodeMap struct {
	hashFunc    HashFunc
	nodeHashs   []int
	nodeHashMap map[int]string // hash值和节点名的map关系
}

func NewNodeMap(fn HashFunc) *NodeMap {
	m := &NodeMap{
		hashFunc:    fn,
		nodeHashMap: make(map[int]string),
	}
	if m.hashFunc == nil {
		m.hashFunc = crc32.ChecksumIEEE
	}
	return m
}

// IsEmpty 判断nodeHashs是否为空
func (m *NodeMap) IsEmpty() bool {
	return len(m.nodeHashs) == 0
}

// AddNode 参数可能是节点ip或者节点名称，必须唯一
func (m *NodeMap) AddNode(keys ...string) {
	for _, key := range keys {
		if key == "" {
			continue
		}
		hash := int(m.hashFunc([]byte(key)))
		m.nodeHashs = append(m.nodeHashs, hash)
		m.nodeHashMap[hash] = key
	}
	// 对切片进行排序
	sort.Ints(m.nodeHashs)
}

// PickNode 根据节点名称，返回
func (m *NodeMap) PickNode(key string) string {
	if m.IsEmpty() {
		return ""
	}
	// 对key进行hash
	hash := int(m.hashFunc([]byte(key)))
	// 获取切片中在那两个节点之间
	// 1000 2000 3000 搜索 1300 得到的idx就是1
	idx := sort.Search(len(m.nodeHashs), func(i int) bool {
		return m.nodeHashs[i] >= hash
	})
	// 判断是否落到最后
	// 如果落到最后，就由第一个节点负责处理
	if idx == len(m.nodeHashs) {
		idx = 0
	}
	// 获取该节点的hash值
	nodeHashVal := m.nodeHashs[idx]
	// 根据hash值获取到实际的ip或者name
	return m.nodeHashMap[nodeHashVal]
}
