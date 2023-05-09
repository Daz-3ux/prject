// 实现一致性哈希
// 		一致性哈希是 GeeCache 从 单节点 走向 多节点 的主要环节
package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

// 依赖注入模式：允许用户替换为自定义 Hash 函数
// 默认为 crc32.ChecksumIEEE 算法
// Hash maps bytes to uint32
type Hash func(data []byte) uint32

// Map contains all hashed keys
type Map struct {
	// 哈希函数
	hash			Hash
	// 虚拟节点倍数
	replicas 	int
	// 哈希环： 避免缓存雪崩
	keys			[]int // Sorted
	// 虚拟节点与真实节点的映射表： 键是虚拟节点的哈希值，值是真实节点的名称
	hashMap		map[int]string
}

// 允许自定义 虚拟节点背书以及 Hash 函数
// New creates a Map instance
func New(replicas int, fn Hash) *Map {
	m := &Map {
		replicas: replicas,
		hash:			fn,
		hashMap: 	make(map[int]string),
	}
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

// 允许传入 0 或多个真实节点名称
// Add adds some keys to the hash
func (m *Map) Add(keys ...string) {
	for _, key := range keys {
		for i := 0; i < m.replicas; i++ {
			// 虚拟节点名称为 strconv.Itoa(i) + key,以编号区分不同虚拟节点
			hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
			// 上环
			m.keys = append(m.keys, hash)
			// 建立映射
			m.hashMap[hash] = key
		}
	}
	// 给环上的哈希值排序
	sort.Ints(m.keys)
}

// Get gets the closest item in the hash to the provided key
func (m *Map) Get(key string) string {
	if len(m.keys) == 0 {
		return ""
	}

	hash := int(m.hash([]byte(key)))
	// 顺时针查找第一个匹配的虚拟节点的下标 idx
	// Binary search for appropriate replica
	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash
	})

	// 环状结构,取余处理
	return m.hashMap[m.keys[idx%len(m.keys)]]
}