package maildoor_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/wawandco/maildoor"
	"github.com/wawandco/maildoor/internal/testhelpers"
)

func TestInMemoryTokenStorage(t *testing.T) {
	t.Run("store and get token", func(t *testing.T) {
		storage := maildoor.NewInMemoryTokenStorage(0) // No expiration

		email := "test@example.com"
		token := "123456"

		err := storage.Store(email, token)
		testhelpers.NoError(t, err)

		retrievedToken, exists := storage.Get(email)
		testhelpers.Equals(t, true, exists)
		testhelpers.Equals(t, token, retrievedToken)
	})

	t.Run("get non-existent token", func(t *testing.T) {
		storage := maildoor.NewInMemoryTokenStorage(0)

		token, exists := storage.Get("nonexistent@example.com")
		testhelpers.Equals(t, false, exists)
		testhelpers.Equals(t, "", token)
	})

	t.Run("delete token", func(t *testing.T) {
		storage := maildoor.NewInMemoryTokenStorage(0)

		email := "test@example.com"
		token := "123456"

		// Store token
		err := storage.Store(email, token)
		testhelpers.NoError(t, err)

		// Verify it exists
		_, exists := storage.Get(email)
		testhelpers.Equals(t, true, exists)

		// Delete token
		deleted := storage.Delete(email)
		testhelpers.Equals(t, true, deleted)

		// Verify it's gone
		_, exists = storage.Get(email)
		testhelpers.Equals(t, false, exists)
	})

	t.Run("delete non-existent token", func(t *testing.T) {
		storage := maildoor.NewInMemoryTokenStorage(0)

		deleted := storage.Delete("nonexistent@example.com")
		testhelpers.Equals(t, false, deleted)
	})

	t.Run("overwrite existing token", func(t *testing.T) {
		storage := maildoor.NewInMemoryTokenStorage(0)

		email := "test@example.com"
		token1 := "123456"
		token2 := "654321"

		// Store first token
		err := storage.Store(email, token1)
		testhelpers.NoError(t, err)

		// Store second token (should overwrite)
		err = storage.Store(email, token2)
		testhelpers.NoError(t, err)

		// Should get the second token
		retrievedToken, exists := storage.Get(email)
		testhelpers.Equals(t, true, exists)
		testhelpers.Equals(t, token2, retrievedToken)
	})

	t.Run("token expiration", func(t *testing.T) {
		storage := maildoor.NewInMemoryTokenStorage(100 * time.Millisecond)

		email := "test@example.com"
		token := "123456"

		// Store token
		err := storage.Store(email, token)
		testhelpers.NoError(t, err)

		// Should exist immediately
		retrievedToken, exists := storage.Get(email)
		testhelpers.Equals(t, true, exists)
		testhelpers.Equals(t, token, retrievedToken)

		// Wait for expiration
		time.Sleep(150 * time.Millisecond)

		// Should be expired
		_, exists = storage.Get(email)
		testhelpers.Equals(t, false, exists)
	})

	t.Run("cleanup expired tokens", func(t *testing.T) {
		storage := maildoor.NewInMemoryTokenStorage(50 * time.Millisecond)

		// Store multiple tokens
		emails := []string{"user1@example.com", "user2@example.com", "user3@example.com"}
		for _, email := range emails {
			err := storage.Store(email, "123456")
			testhelpers.NoError(t, err)
		}

		// All should exist
		for _, email := range emails {
			_, exists := storage.Get(email)
			testhelpers.Equals(t, true, exists)
		}

		// Wait for expiration
		time.Sleep(100 * time.Millisecond)

		// Manually trigger cleanup
		storage.Cleanup()

		// All should be gone after cleanup
		for _, email := range emails {
			_, exists := storage.Get(email)
			testhelpers.Equals(t, false, exists)
		}
	})

	t.Run("no expiration means no cleanup", func(t *testing.T) {
		storage := maildoor.NewInMemoryTokenStorage(0) // No expiration

		email := "test@example.com"
		token := "123456"

		err := storage.Store(email, token)
		testhelpers.NoError(t, err)

		// Wait a bit
		time.Sleep(50 * time.Millisecond)

		// Trigger cleanup (should do nothing)
		storage.Cleanup()

		// Token should still exist
		retrievedToken, exists := storage.Get(email)
		testhelpers.Equals(t, true, exists)
		testhelpers.Equals(t, token, retrievedToken)
	})
}

func TestMaildoorWithCustomTokenStorage(t *testing.T) {
	t.Run("custom token storage integration", func(t *testing.T) {
		// Create custom storage with expiration
		storage := maildoor.NewInMemoryTokenStorage(5 * time.Minute)

		var sentEmails []string
		var sentTokens []string

		auth := maildoor.New(
			maildoor.WithTokenStorage(storage),
			maildoor.EmailValidator(func(email string) error {
				return nil
			}),
			maildoor.EmailSender(func(email, html, txt string) error {
				sentEmails = append(sentEmails, email)
				// Extract token from message for testing
				lines := strings.Split(txt, "\n")
				for _, line := range lines {
					if strings.Contains(line, "Code:") {
						parts := strings.Split(line, "Code:")
						if len(parts) > 1 {
							token := strings.TrimSpace(parts[1])
							sentTokens = append(sentTokens, token)
						}
					}
				}
				return nil
			}),
		)

		// Request token for email
		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/email", nil)
		req.Form = url.Values{
			"email": []string{"test@example.com"},
		}

		auth.ServeHTTP(w, req)
		testhelpers.Equals(t, http.StatusOK, w.Code)

		// Verify email was sent
		testhelpers.Equals(t, 1, len(sentEmails))
		testhelpers.Equals(t, "test@example.com", sentEmails[0])
		testhelpers.Equals(t, 1, len(sentTokens))

		// Verify token is stored in custom storage
		storedToken, exists := storage.Get("test@example.com")
		testhelpers.Equals(t, true, exists)
		testhelpers.Equals(t, sentTokens[0], storedToken)

		// Test code validation
		w = httptest.NewRecorder()
		req = httptest.NewRequest("POST", "/code", nil)
		req.Form = url.Values{
			"email": []string{"test@example.com"},
			"code":  []string{sentTokens[0]},
		}

		auth.ServeHTTP(w, req)
		testhelpers.Equals(t, http.StatusOK, w.Code)

		// Verify token was deleted after successful login
		_, exists = storage.Get("test@example.com")
		testhelpers.Equals(t, false, exists)
	})
}

// MockTokenStorage is a test implementation of TokenStorage
type MockTokenStorage struct {
	tokens      map[string]string
	storeError  error
	deleteError bool
}

func NewMockTokenStorage() *MockTokenStorage {
	return &MockTokenStorage{
		tokens: make(map[string]string),
	}
}

func (m *MockTokenStorage) Store(email, token string) error {
	if m.storeError != nil {
		return m.storeError
	}
	m.tokens[email] = token
	return nil
}

func (m *MockTokenStorage) Get(email string) (string, bool) {
	token, exists := m.tokens[email]
	return token, exists
}

func (m *MockTokenStorage) Delete(email string) bool {
	if m.deleteError {
		return false
	}
	_, exists := m.tokens[email]
	if exists {
		delete(m.tokens, email)
	}
	return exists
}

func (m *MockTokenStorage) Cleanup() {
	// Mock implementation - do nothing
}

func TestMaildoorWithMockTokenStorage(t *testing.T) {
	t.Run("mock token storage", func(t *testing.T) {
		mockStorage := NewMockTokenStorage()

		auth := maildoor.New(
			maildoor.WithTokenStorage(mockStorage),
			maildoor.EmailValidator(func(email string) error {
				return nil
			}),
			maildoor.EmailSender(func(email, html, txt string) error {
				return nil
			}),
		)

		// Request token
		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/email", nil)
		req.Form = url.Values{
			"email": []string{"test@example.com"},
		}

		auth.ServeHTTP(w, req)
		testhelpers.Equals(t, http.StatusOK, w.Code)

		// Verify token was stored in mock storage
		token, exists := mockStorage.Get("test@example.com")
		testhelpers.Equals(t, true, exists)
		testhelpers.Equals(t, 6, len(token)) // Should be 6 digits

		// Test successful code validation
		w = httptest.NewRecorder()
		req = httptest.NewRequest("POST", "/code", nil)
		req.Form = url.Values{
			"email": []string{"test@example.com"},
			"code":  []string{token},
		}

		auth.ServeHTTP(w, req)
		testhelpers.Equals(t, http.StatusOK, w.Code)

		// Token should be deleted after successful login
		_, exists = mockStorage.Get("test@example.com")
		testhelpers.Equals(t, false, exists)
	})
}
