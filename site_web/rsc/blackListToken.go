package API

import (
	"container/list"
	"sync"
	"time"
)

var blacklist *TokenBlacklist

// TokenBlacklist represents a token blacklist with a maximum size and FIFO eviction policy
type TokenBlacklist struct {
	maxSize int
	list    *list.List
	mutex   sync.Mutex
}

// BlacklistedToken represents a blacklisted token
type BlacklistedToken struct {
	Token      string
	ExpiryTime time.Time
}

func InitBlacklist() {
	blacklist = NewTokenBlacklist(100)
}

// NewTokenBlacklist creates a new TokenBlacklist with the specified maximum size
func NewTokenBlacklist(maxSize int) *TokenBlacklist {
	return &TokenBlacklist{
		maxSize: maxSize,
		list:    list.New(),
	}
}

// AddToken adds a token to the blacklist
func (tb *TokenBlacklist) AddToken(token string, expiryTime time.Time) {
	tb.mutex.Lock()
	defer tb.mutex.Unlock()

	// Evict oldest token if blacklist is full
	if tb.list.Len() >= tb.maxSize {
		tb.list.Remove(tb.list.Front())
	}

	// Add new token to the end of the list
	tb.list.PushBack(BlacklistedToken{
		Token:      token,
		ExpiryTime: expiryTime,
	})
}

// IsTokenBlacklisted checks if a token is blacklisted
func (tb *TokenBlacklist) IsTokenBlacklisted(token string) bool {
	tb.mutex.Lock()
	defer tb.mutex.Unlock()

	// Check if the token is in the blacklist
	for elem := tb.list.Front(); elem != nil; elem = elem.Next() {
		blacklistedToken := elem.Value.(BlacklistedToken)
		if blacklistedToken.Token == token {
			// Remove expired tokens
			if blacklistedToken.ExpiryTime.Before(time.Now()) {
				tb.list.Remove(elem)
				return false
			}
			return true
		}
	}
	return false
}
