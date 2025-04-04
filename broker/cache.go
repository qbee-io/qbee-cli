package broker

import (
	"sync"
	"time"
)

type item struct {
	data    any
	expires time.Time
}

type Cache struct {
	items map[string]item
	ttl   time.Duration
	mutex sync.Mutex
}

func (c *Cache) Add(key string, data any) {

	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.items[key] = item{
		data:    data,
		expires: time.Now().Add(c.ttl),
	}
}

func (c *Cache) Get(key string) (any, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	item, ok := c.items[key]
	if !ok {
		return nil, false
	}
	return item.data, ok
}

func (c *Cache) Remove(key string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	delete(c.items, key)
}

func NewCache(ttl time.Duration) Cache {
	return Cache{
		items: make(map[string]item),
		ttl:   ttl,
	}
}
