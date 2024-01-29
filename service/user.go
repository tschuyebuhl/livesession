package service

import (
	"github.com/tschuyebuhl/livesession/cache"
	"github.com/tschuyebuhl/livesession/data"
	"log/slog"
	"sync"
)

type UserService struct {
	repo    data.Repository
	cache   cache.Cache
	pending map[data.ID][]chan *data.User
	mu      sync.Mutex
}

func NewUserService(repo data.Repository, cache cache.Cache) *UserService {
	return &UserService{
		repo:    repo,
		cache:   cache,
		pending: make(map[data.ID][]chan *data.User),
	}
}

// GetUser returns a user from the cache if it exists, otherwise it fetches it from the repository and caches it
// double-checked locking is used in case a goroutine has updated the cache
func (s *UserService) GetUser(id data.ID) (*data.User, bool, error) {
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
		response := make(chan *data.User)
		s.pending[id] = append(waiters, response)
		s.mu.Unlock()
		return <-response, true, nil
	}

	slog.Info("no pending requests, marking as pending", "id", id)
	s.pending[id] = []chan *data.User{}
	s.mu.Unlock()

	slog.Info("cache missed, fetching from repo", "id", id)
	user, err := s.repo.Get(id)
	if err != nil {
		s.notifyWaiters(nil, id) // notify with nil when query fails
		return nil, false, err
	}

	s.cache.Put(user)
	s.notifyWaiters(user, id)
	return user, false, nil
}

func (s *UserService) notifyWaiters(user *data.User, id data.ID) {
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
