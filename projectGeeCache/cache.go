/*
		使用 sync.Mutex 封装 LRU 方法
		使其支持并发的读写
*/
package geecache

import (
	"lru"
	"sync"
)

type cache struct {
	mu					sync.Mutex
	lru					*lru.Cache
	cacheBytes	int64
}

func (c *cache) add(key string, value ByteView) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		c.lru = lru.New(c.cacheBytes, nil)
	}
	c.lru.Add(key, value)
}

func (c *cache) get(key string) (value ByteView, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		return
	}

	if v, ok := c.lru.Get(key); ok {
		return v.(ByteView), ok
	}

	return
}

