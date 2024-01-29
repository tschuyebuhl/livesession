package cache

import (
	"github.com/tschuyebuhl/livesession/data"
	"sync"
)

type InMem struct {
	data map[data.ID]*data.User
	mu   sync.RWMutex
}

func NewInMemoryCache() *InMem {
	return &InMem{
		data: make(map[data.ID]*data.User),
	}
}

func (c *InMem) Delete(key data.ID) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.data, key)
}

func (c *InMem) Nuke(sure bool) {
	if sure {
		c.mu.Lock()
		defer c.mu.Unlock()
		c.data = make(map[data.ID]*data.User)
	}
}

func (c *InMem) Put(val *data.User) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[val.ID] = val
}

func (c *InMem) Get(key data.ID) (*data.User, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	val, ok := c.data[key]
	return val, ok
}
