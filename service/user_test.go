package service

import (
	"fmt"
	"github.com/tschuyebuhl/livesession/cache"
	"github.com/tschuyebuhl/livesession/data"
	"sync"
	"testing"
)

type MockRepository struct {
	mu         sync.Mutex // To protect access to the counter
	queryCount int
}

func (m *MockRepository) Get(id data.ID) (*data.User, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.queryCount++
	// Simulate database access
	return &data.User{ID: id, Name: "Test", Surname: "User"}, nil
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

func (m *MockRepository) QueryCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.queryCount
}

// TestUserServiceLoad tests the cache mechanism under load
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
