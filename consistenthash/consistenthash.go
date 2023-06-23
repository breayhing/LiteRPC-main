package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

type ConsistentHash struct {
	replica int   // 虚拟节点倍数
	keys    []int // 哈希环
	hashMap map[int]string
}

func NewConsistentHash(replica int) *ConsistentHash {
	return &ConsistentHash{
		replica: replica,
		hashMap: make(map[int]string),
	}
}

// 增加节点
func (c *ConsistentHash) Add(key string) {
	for i := 0; i < c.replica; i++ {
		hash := int(crc32.ChecksumIEEE([]byte(strconv.Itoa(i) + key)))
		c.keys = append(c.keys, hash)
		c.hashMap[hash] = key
	}
	sort.Ints(c.keys)
}

// 获取key对应的节点
func (c *ConsistentHash) Get(key string) string {
	if len(c.keys) == 0 {
		return ""
	}
	hash := int(crc32.ChecksumIEEE([]byte(key)))
	idx := sort.SearchInts(c.keys, hash)
	return c.hashMap[c.keys[idx%len(c.keys)]]
}
