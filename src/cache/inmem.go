package cache

import (
	"github.com/tschuyebuhl/livesession/src/data"
	"log/slog"
	"sync"
)

type InMem struct {
	data map[data.ID]*data.User
	mu   sync.RWMutex
}

func NewInMemory() *InMem {
	return &InMem{
		data: make(map[data.ID]*data.User),
	}
}

func (c *InMem) Delete(key data.ID) {
	slog.Info("deleting data from cache", "id", key)
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.data, key)
}

func (c *InMem) Nuke(sure bool) {
	slog.Info("nuking cache")
	if sure {
		c.mu.Lock()
		defer c.mu.Unlock()
		c.data = make(map[data.ID]*data.User)
	}
}

func (c *InMem) Put(val *data.User) {
	slog.Info("putting data in cache", "id", val.ID)
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[val.ID] = val
}

func (c *InMem) Get(key data.ID) (*data.User, bool, bool) {
	slog.Info("getting data from cache", "id", key)
	c.mu.RLock()
	if val, ok := c.data[key]; ok {
		c.mu.RUnlock()
		return val, true, true
	}
	c.mu.RUnlock()
	return nil, false, false
}
