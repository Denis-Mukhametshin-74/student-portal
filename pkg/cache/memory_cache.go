package cache

import (
	"sync"
	"time"
)

type CacheItem struct {
	Value      any
	Expiration int64
}

type MemoryCache struct {
	items map[string]CacheItem
	mu    sync.RWMutex
}

func NewMemoryCache() *MemoryCache {
	cache := &MemoryCache{
		items: make(map[string]CacheItem),
	}

	go cache.startCleanup()

	return cache
}

func (c *MemoryCache) Set(key string, value any, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = CacheItem{
		Value:      value,
		Expiration: time.Now().Add(ttl).UnixNano(),
	}
}

func (c *MemoryCache) Get(key string) (any, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, exists := c.items[key]
	if !exists {
		return nil, false
	}

	if time.Now().UnixNano() > item.Expiration {
		return nil, false
	}

	return item.Value, true
}

func (c *MemoryCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
}

func (c *MemoryCache) startCleanup() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		c.mu.Lock()
		now := time.Now().UnixNano()

		for key, item := range c.items {
			if now > item.Expiration {
				delete(c.items, key)
			}
		}
		c.mu.Unlock()
	}
}
