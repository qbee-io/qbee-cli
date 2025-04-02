package broker

import (
	"context"
	"sync"
	"time"
)

type Connection struct {
	Cancel        context.CancelFunc
	Port          int
	LastConnected time.Time
}

type ConnectionsCache struct {
	connections map[string]Connection
	mutex       sync.Mutex
}

func (c *ConnectionsCache) Add(key string, conn Connection) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.connections[key] = conn
}

func (c *ConnectionsCache) Get(key string) (Connection, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	conn, ok := c.connections[key]
	return conn, ok
}

func (c *ConnectionsCache) Remove(key string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	delete(c.connections, key)
}

func (c *ConnectionsCache) CleanUp() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	for k, v := range c.connections {
		if time.Since(v.LastConnected) > defaultConnectionTTL {
			v.Cancel()
			delete(c.connections, k)
		}
	}
}

func NewConnectionsCache() ConnectionsCache {
	return ConnectionsCache{
		connections: make(map[string]Connection),
	}
}
