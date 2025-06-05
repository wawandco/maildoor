package maildoor

import (
	"sync"
	"time"
)

// TokenStorage defines the interface for storing and retrieving authentication tokens.
// This allows for custom storage implementations such as Redis, database, or other backends.
type TokenStorage interface {
	// Store saves a token for the given email address.
	// If a token already exists for the email, it should be overwritten.
	Store(email, token string) error

	// Get retrieves the token for the given email address.
	// Returns the token and true if found, empty string and false if not found.
	Get(email string) (token string, exists bool)

	// Delete removes the token for the given email address.
	// Returns true if the token existed and was deleted, false if it didn't exist.
	Delete(email string) bool
}

// InMemoryTokenStorage is the default in-memory implementation of ITokenStorage.
// It stores tokens in memory with optional expiration support.
type InMemoryTokenStorage struct {
	mu     sync.RWMutex
	tokens map[string]tokenEntry
	ttl    time.Duration
}

type tokenEntry struct {
	token     string
	createdAt time.Time
}

// NewInMemoryTokenStorage creates a new in-memory token storage.
// If ttl is 0, tokens never expire. If ttl > 0, tokens expire after the specified duration.
func NewInMemoryTokenStorage(ttl time.Duration) *InMemoryTokenStorage {
	storage := &InMemoryTokenStorage{
		tokens: make(map[string]tokenEntry),
		ttl:    ttl,
	}

	// Start cleanup goroutine if TTL is set
	if ttl > 0 {
		go storage.cleanupLoop()
	}

	return storage
}

// Store implements ITokenStorage.Store
func (s *InMemoryTokenStorage) Store(email, token string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.tokens[email] = tokenEntry{
		token:     token,
		createdAt: time.Now(),
	}

	return nil
}

// Get implements ITokenStorage.Get
func (s *InMemoryTokenStorage) Get(email string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	entry, exists := s.tokens[email]
	if !exists {
		return "", false
	}

	// Check if token has expired
	if s.ttl > 0 && time.Since(entry.createdAt) > s.ttl {
		// Remove expired token
		s.mu.RUnlock()
		s.mu.Lock()
		delete(s.tokens, email)
		s.mu.Unlock()
		s.mu.RLock()
		return "", false
	}

	return entry.token, true
}

// Delete implements ITokenStorage.Delete
func (s *InMemoryTokenStorage) Delete(email string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, exists := s.tokens[email]
	if exists {
		delete(s.tokens, email)
	}

	return exists
}

// Cleanup implements ITokenStorage.Cleanup
func (s *InMemoryTokenStorage) Cleanup() {
	if s.ttl == 0 {
		return // No expiration, nothing to cleanup
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	for email, entry := range s.tokens {
		if now.Sub(entry.createdAt) > s.ttl {
			delete(s.tokens, email)
		}
	}
}

// cleanupLoop runs periodic cleanup of expired tokens
func (s *InMemoryTokenStorage) cleanupLoop() {
	// Cleanup every half of the TTL duration, with a minimum of 1 minute
	interval := s.ttl / 2
	if interval < time.Minute {
		interval = time.Minute
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		s.Cleanup()
	}
}
