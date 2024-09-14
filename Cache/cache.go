package Cache

import (
	"7days/Cache/lru"
	"sync"
)

//控制并发,key是存储元素用的键，name是缓存空间的名字
//真正的缓存，lru只是一种存储方式，后面可以考虑将消息种类分类，用不同算法存储

type cache struct {
	maxcache int64 //最大容量
	lru      *lru.Cache
	mu       sync.Mutex
}

func (c *cache) add(key string, value ByteView) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		c.lru = lru.New(c.maxcache, nil)
	}
	c.lru.Add(key, value)
}

func (c *cache) get(key string) (value ByteView, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if key == "" {
		return
	}
	if c.lru == nil {
		return
	}
	if v, ok := c.lru.Get(key); ok {
		return v.(ByteView), ok
	}
	return
}
