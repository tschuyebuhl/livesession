package service

import (
	"github.com/tschuyebuhl/livesession/src/cache"
	data2 "github.com/tschuyebuhl/livesession/src/data"
	"log/slog"
	"sync"
)

type UserService struct {
	repo    data2.Repository
	cache   cache.Cache
	pending map[data2.ID][]chan *data2.User
	mu      sync.Mutex
}

func NewUserService(repo data2.Repository, cache cache.Cache) *UserService {
	return &UserService{
		repo:    repo,
		cache:   cache,
		pending: make(map[data2.ID][]chan *data2.User),
	}
}

// GetUser returns a user from the cache if it exists, otherwise it fetches it from the repository and caches it
// double-checked locking is used in case a goroutine has updated the cache in the meantime
func (s *UserService) GetUser(id data2.ID) (*data2.User, bool, error) {
	slog.Info("getting user", "id", id)
	if user, found, _ := s.cache.Get(id); found {
		return user, true, nil
	}

	s.mu.Lock()
	if user, found, _ := s.cache.Get(id); found {
		s.mu.Unlock()
		return user, true, nil
	}
	if waiters, ok := s.pending[id]; ok {
		response := make(chan *data2.User)
		s.pending[id] = append(waiters, response)
		s.mu.Unlock()
		return <-response, true, nil
	}

	slog.Info("no pending requests, marking as pending", "id", id)
	s.pending[id] = []chan *data2.User{}
	s.mu.Unlock()

	slog.Info("cache missed, fetching from repo", "id", id)
	user, err := s.repo.Get(id)
	if err != nil {
		slog.Error("error fetching from repo", "id", id, "error", err)
		s.notifyWaiters(nil, id) // notify with nil when query fails
		return nil, false, err
	}

	s.cache.Put(user)
	s.notifyWaiters(user, id)
	return user, false, nil
}

func (s *UserService) notifyWaiters(user *data2.User, id data2.ID) {
	s.mu.Lock()
	defer s.mu.Unlock()

	slog.Info("notifying waiters", "id", id)
	if waiters, ok := s.pending[id]; ok {
		for _, waiter := range waiters {
			waiter <- user
		}
		delete(s.pending, id)
	}
}
