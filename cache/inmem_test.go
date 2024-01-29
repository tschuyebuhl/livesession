package cache

import (
	"github.com/tschuyebuhl/livesession/data"
	"testing"
)

func TestCacheOperations(t *testing.T) {
	cache := NewInMemory()

	// Mock user data
	user1 := &data.User{ID: "user1"}

	cache.Put(user1)
	if got, found, _ := cache.Get("user1"); !found || got.ID != user1.ID {
		t.Errorf("Get(user1) = %v, want %v", got, user1)
	}

	if _, found, cached := cache.Get("user2"); found && !cached {
		t.Errorf("Expected cache miss for user2 but found cached data")
	}

}
