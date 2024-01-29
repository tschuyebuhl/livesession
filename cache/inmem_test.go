package cache

import (
	"github.com/tschuyebuhl/livesession/data"
	"sync"
	"testing"
)

func TestCacheOperations(t *testing.T) {
	cache := NewInMemoryCache()
	finished := make(chan bool)
	var wg sync.WaitGroup

	// Mock user data
	user1 := &data.User{ID: "user1"}
	user2 := &data.User{ID: "user2"}

	cache.Put(user1)
	if got, found, _ := cache.Get("user1"); !found || got.ID != user1.ID {
		t.Errorf("Get(user1) = %v, want %v", got, user1)
	}

	if _, found, cached := cache.Get("user2"); found && cached {
		t.Errorf("Expected cache miss for user2 but found cached data")
	}

	// now the user should be cached
	if got, found, _ := cache.Get("user2"); !found || got.ID != user2.ID {
		t.Errorf("Get(user2) = %v, want %v", got, user2)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		cache.Put(&data.User{ID: "user3"})
		finished <- true
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		<-finished
		if _, found, _ := cache.Get("user3"); !found {
			t.Errorf("Concurrent access failed for user3")
		}
	}()

	wg.Wait()
}
