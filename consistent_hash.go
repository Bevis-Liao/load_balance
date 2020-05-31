package load_balance

import (
	"errors"
	"hash/crc32"
	"sort"
	"strconv"
	"sync"
)

type Hash func(data []byte) uint32

type UInt32Slice []uint32

func (s UInt32Slice) Len() int {
	return len(s)
}

func (s UInt32Slice) Less(i, j int) bool {
	return s[i] < s[j]
}

func (s UInt32Slice) Swap(i, j int)  {
	s[i], s[j] = s[j], s[i]
}

type ConsistentHashBalance struct {
	mux sync.RWMutex
	hash Hash
	replicas int // 复制因子
	keys UInt32Slice // 已排序的节点 hash 切片
	hasMap map[uint32]string // 节点哈希和 key 的 map, 键是 hash 值，值是节点 key
}

func NewConsistentHashBalance(replicas int, fn Hash) *ConsistentHashBalance {
	m := &ConsistentHashBalance{
		replicas: replicas,
		hash:     fn,
		hasMap:   make(map[uint32]string),
	}
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}

	return m
}

func (c *ConsistentHashBalance) IsEmpty() bool {
	return len(c.keys) == 0
}

func (c *ConsistentHashBalance) Add(params ...string) error {
	if len(params) == 0 {
		return errors.New("param len 1 at least")
	}
	addr := params[0]
	c.mux.Lock()
	defer c.mux.Unlock()
	for i := 0; i < c.replicas; i++ {
		hash := c.hash([]byte(strconv.Itoa(i) + addr))
		c.keys = append(c.keys, hash)
		c.hasMap[hash] = addr
	}
	// 对所有虚拟节点的哈希值进行排序，方便之和进行二分查找
	sort.Sort(c.keys)
	return nil
}

func (c *ConsistentHashBalance) Get(key string) (string, error)  {
	if c.IsEmpty() {
		return "", errors.New("node is empty")
	}

	hash := c.hash([]byte(key))

	// 通过二分查找获取最优节点，第一个"服务器 hash" 值大于 "数据 hash" 值的就是最优"服务器节点"
	idx := sort.Search(len(c.keys), func(i int) bool {
		return c.keys[i] >= hash
	})

	// 如果查找结果大于服务器节点的哈希数组的最大索引，表示此时该对象哈希值位于一个节点之后，那么放
	if idx == len(c.keys) {
		idx = 0
	}

	c.mux.RLock()
	defer c.mux.RUnlock()
	return c.hasMap[c.keys[idx]], nil
}

// 一致性首先计算四个节点对应的 hash 值
// hash 环上顺时针从整数 0 开始，一直到最大正整数，ip hash 值肯定会落到这个 Hash 环上的某一点，即把 ip 映射到环上。
// 当用户请求时，首先会 hash 用户的 ip，然后看落在 hash 环的哪个地方，根据 hash 环上的位置 顺时针 找距离最近的 ip 作为路由 ip