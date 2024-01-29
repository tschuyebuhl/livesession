package data

import "github.com/tschuyebuhl/livesession/cache"

// ID is used as a key, so it must be comparable, type alias is used to allow plugging in different types
type ID string

type User struct {
	ID      ID
	Name    string
	Surname string
}

type UserService struct {
	repo  Repository
	cache cache.Cache
}

func NewUserService(repo Repository, cache cache.Cache) *UserService {
	return &UserService{
		repo:  repo,
		cache: cache,
	}
}

func (s *UserService) FetchUser(id ID) (*User, error) {
	return s.repo.Get(id)
}
