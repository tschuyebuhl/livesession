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

func (s *UserService) GetUser(id data.ID) (*data.User, bool, error) {
	slog.Info("getting user", "id", id)
	if user, found, _ := s.cache.Get(id); found {
		return user, true, nil
	}

	s.mu.Lock()
	if waiters, ok := s.pending[id]; ok {
		response := make(chan *data.User)
		s.pending[id] = append(waiters, response)
		s.mu.Unlock()
		user := <-response
		return user, user != nil, nil
	}

	response := make(chan *data.User)
	s.pending[id] = []chan *data.User{response}
	s.mu.Unlock()

	user, err := s.repo.Get(id)
	if err != nil {
		s.mu.Lock()
		delete(s.pending, id)
		s.notifyWaiters(nil, id)
		s.mu.Unlock()
		return nil, false, err
	}

	s.cache.Put(user)
	s.notifyWaiters(user, id)
	return user, false, nil
}

func (s *UserService) notifyWaiters(user *data.User, id data.ID) {
	slog.Info("notifying waiters", "id", id)
	if waiters, ok := s.pending[id]; ok {
		for _, waiter := range waiters {
			waiter <- user
		}
		delete(s.pending, id)
	}
}
