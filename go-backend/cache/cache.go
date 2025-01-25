package cache

import (
	"sync"
	"time"
)

var (
	// Stores the user credentials and data
	userCache = sync.Map{}
)

// Caches user data for 24 hours
func CacheUser(userID string, credentials interface{}) {
	// userID is key, credentials is value
	userCache.Store(userID, credentials)
	time.AfterFunc(24 * time.Hour, func() {
		userCache.Delete(userID)
	})
}

// Retrieves user credential if token is not expired
func GetCachedUser(userID string) (interface{}, bool) {
	return userCache.Load(userID)
}

func DeleteCachedUser(userID string) {
	userCache.Delete(userID)
}