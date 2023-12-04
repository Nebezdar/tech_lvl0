package cache

import (
	"sync"
	"time"
)

type Cache struct {
	data      map[string][]byte
	expiry    int
	mutex     sync.Mutex
	cleanupWg sync.WaitGroup
}

func NewCache(expiry int) *Cache {
	return &Cache{
		data:   make(map[string][]byte),
		expiry: expiry,
	}
}

func (c *Cache) Set(key string, value []byte) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.data[key] = value

	// Remove expired entries from the cache
	go func() {
		time.Sleep(time.Duration(c.expiry) * time.Second)
		c.mutex.Lock()
		defer c.mutex.Unlock()
		delete(c.data, key)
	}()
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	value, ok := c.data[key]
	return value, ok
}
