// lru 缓存淘汰策略
package lru

import (
	"container/list"
)

// Cache is a LRU cache, not safe fot concurrent
type Cache struct {
	// maxBytes 为 0 时代表不对内存大小设限
	maxBytes int64
	// nBytes 为已使用内存
	nBytes   int64
	// 双向链表: Front 放置最新使用节点
	ll       *list.List
	// 缓存
	cache    map[string]*list.Element
	// 回调函数: optional and executed when an entry is purged
	onEvicated func(key string, value Value)
}

// 双向链表的数据类型,定义 Value any, 定义接口
type entry struct {
	key   string
	value Value
}

// Value use len to count how many bytes it takes
type Value interface {
	Len() int
}

// New is the Constructor of Cache
func New(maxBytes int64, onEvicated func(string, Value)) *Cache {
	return &Cache{
		maxBytes:   maxBytes,
		ll:         list.New(),
		cache:      make(map[string]*list.Element),
		onEvicated: onEvicated,
	}
}

// Get look ups a key's value
func (c *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := c.cache[key]; ok {
		// Conventionally, front refers to the tail of the queue -- 约定 front 作为队尾
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	return
}

// RemoveOldest removes the oldest item
func (c *Cache) RemoveOldest() {
	// 从队首取元素
	ele := c.ll.Back()
	if ele != nil {
		c.ll.Remove(ele)
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)
		c.nBytes -= int64(len(kv.key)) + int64(kv.value.Len())
		if c.onEvicated != nil {
			c.onEvicated(kv.key, kv.value)
		}
	}
}

// Add adds a value to the cache
func (c *Cache) Add(key string, value Value) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		c.nBytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else {
		ele := c.ll.PushFront(&entry{key, value})
		c.cache[key] = ele
		c.nBytes += int64(len(key)) + int64(value.Len())
	}
	for c.maxBytes != 0 && c.maxBytes < c.nBytes {
		c.RemoveOldest()
	}
}

// Len the number of cache entries
func (c *Cache) Len() int {
	return c.ll.Len()
}
