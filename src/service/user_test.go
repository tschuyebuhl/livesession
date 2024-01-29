package service

import (
	"fmt"
	"github.com/tschuyebuhl/livesession/src/cache"
	"github.com/tschuyebuhl/livesession/src/data"
	"sync"
	"testing"
)

type MockRepository struct {
	mu         sync.Mutex // to protect counter access
	queryCount int
	users      map[data.ID]*data.User
}

func (m *MockRepository) Get(id data.ID) (*data.User, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.queryCount++
	return &data.User{ID: id, Name: "Test", Surname: "User"}, nil
}

func (m *MockRepository) QueryCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.queryCount
}

func (m *MockRepository) SetUser(user *data.User) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.users == nil {
		m.users = make(map[data.ID]*data.User)
	}
	m.users[user.ID] = user
}

// tests for race conditions
func TestUserServiceConcurrentAccess(t *testing.T) {
	repo := &MockRepository{}
	c := cache.NewInMemory()
	userService := NewUserService(repo, c)

	var wg sync.WaitGroup
	numConcurrentRequests := 10

	userID := data.ID("testUserID")

	for i := 0; i < numConcurrentRequests; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, _, err := userService.GetUser(userID)
			if err != nil {
				t.Errorf("GetUser returned an error: %v", err)
			}
		}()
	}
	wg.Wait()
}

// TestUserServiceLoad tests that the cache is used when multiple goroutines are requesting the same user
// and that the repository is only queried once per id
func TestUserServiceLoad(t *testing.T) {
	repo := &MockRepository{}
	c := cache.NewInMemory()
	userService := NewUserService(repo, c)

	var wg sync.WaitGroup
	numRequests := 10000
	numIDs := 100

	for i := 0; i < numRequests; i++ {
		wg.Add(1)
		userID := data.ID(fmt.Sprintf("userID%d", i%numIDs))

		go func(id data.ID) {
			defer wg.Done()
			_, _, err := userService.GetUser(id)
			if err != nil {
				t.Errorf("GetUser returned an error: %v", err)
			}
		}(userID)
	}
	wg.Wait()

	if repo.QueryCount() != numIDs {
		t.Errorf("Expected %d queries, but got %d", numIDs, repo.QueryCount())
	}
}

func TestUserServiceCorrectUserData(t *testing.T) {
	repo := &MockRepository{}
	c := cache.NewInMemory()
	userService := NewUserService(repo, c)
	repo.SetUser(&data.User{ID: "testUserID", Name: "Test", Surname: "User"})

	userID := data.ID("testUserID")
	user, _, err := userService.GetUser(userID)
	if err != nil {
		t.Errorf("GetUser returned an error: %v", err)
	}
	if user.ID != userID || user.Name != "Test" || user.Surname != "User" {
		t.Errorf("GetUser returned incorrect user data: got %v, want ID: %s, Name: Test, Surname: User", user, userID)
	}
}

func TestCacheAndRepositoryInteraction(t *testing.T) {
	repo := &MockRepository{}
	c := cache.NewInMemory()
	userService := NewUserService(repo, c)
	repo.SetUser(&data.User{ID: "testUserID", Name: "Test", Surname: "User"})
	userID := data.ID("testUserID")

	// should be a cache miss
	_, _, err := userService.GetUser(userID)
	if err != nil {
		t.Errorf("GetUser returned an error on first call: %v", err)
	}

	// should be a cache hit
	_, _, err = userService.GetUser(userID)
	if err != nil {
		t.Errorf("GetUser returned an error on second call: %v", err)
	}

	if repo.QueryCount() != 1 {
		t.Errorf("Expected 1 query to repository, but got %d", repo.QueryCount())
	}
}
