package cache

import (
	"github.com/tschuyebuhl/livesession/data"
	"log/slog"
	"sync"
)

type InMem struct {
	data    map[data.ID]*data.User
	pending map[data.ID][]chan *data.User
	mu      sync.RWMutex
}

func NewInMemoryCache() *InMem {
	return &InMem{
		data:    make(map[data.ID]*data.User),
		pending: make(map[data.ID][]chan *data.User),
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

func (c *InMem) Get(key data.ID) (*data.User, bool, bool) {
	slog.Info("getting user", "key", key)
	slog.Debug("locking cache for r/w")
	c.mu.RLock()
	if val, ok := c.data[key]; ok {
		slog.Info("found user in cache", "key", key)
		c.mu.RUnlock()
		return val, true, true
	}

	slog.Debug("checking pending requests")
	if waiters, ok := c.pending[key]; ok {
		response := make(chan *data.User)
		c.pending[key] = append(waiters, response)
		c.mu.RUnlock()
		return <-response, true, true
	}

	slog.Debug("no pending requests, creating one")
	c.pending[key] = []chan *data.User{}
	c.mu.RUnlock()

	slog.Warn("cache missed, fetching user") // warn may not be appropriate here
	user, err := data.FetchUser(key)
	if err != nil {
		delete(c.pending, key)
		return nil, false, false
	}

	slog.Info("fetch ok, caching user", "key", key)
	c.Put(user)

	c.mu.Lock()
	slog.Debug("notifying all waiters")
	for _, waiter := range c.pending[key] {
		waiter <- user
	}

	delete(c.pending, key)
	slog.Debug("deleted pending requests")
	c.mu.Unlock()

	return user, true, false
}
