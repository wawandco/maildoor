package maildoor

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestCustomLoginRenderer(t *testing.T) {
	customHTML := "<html><body>Custom Login Page</body></html>"
	
	handler := New(
		LoginRenderer(func(data Attempt) (string, error) {
			if data.ProductName != "Maildoor" {
				t.Errorf("Expected ProductName to be 'Maildoor', got %s", data.ProductName)
			}
			return customHTML, nil
		}),
	)

	req := httptest.NewRequest("GET", "/login", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", w.Code)
	}

	if w.Body.String() != customHTML {
		t.Errorf("Expected custom HTML, got %s", w.Body.String())
	}

	if w.Header().Get("Content-Type") != "text/html" {
		t.Errorf("Expected Content-Type to be text/html, got %s", w.Header().Get("Content-Type"))
	}
}

func TestCustomCodeRenderer(t *testing.T) {
	customHTML := "<html><body>Custom Code Page</body></html>"
	testEmail := "test@example.com"
	
	// Create a custom token storage for testing
	tokenStorage := NewInMemoryTokenStorage(0)
	
	handler := New(
		TokenStorage(tokenStorage),
		CodeRenderer(func(data Attempt) (string, error) {
			if data.Email != testEmail {
				t.Errorf("Expected Email to be %s, got %s", testEmail, data.Email)
			}
			if data.Error != "Invalid token" {
				t.Errorf("Expected Error to be 'Invalid token', got %s", data.Error)
			}
			return customHTML, nil
		}),
		EmailSender(func(to, html, txt string) error {
			return nil
		}),
	)

	// Set a specific code for the test email to guarantee we know what it is
	err := tokenStorage.Store(testEmail, "123456")
	if err != nil {
		t.Errorf("Failed to store token: %v", err)
	}

	// Submit an invalid code to trigger error and custom renderer
	form := url.Values{}
	form.Add("email", testEmail)
	form.Add("code", "999999")  // This is guaranteed to be different from "123456"
	
	req := httptest.NewRequest("POST", "/code", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", w.Code)
	}

	if w.Body.String() != customHTML {
		t.Errorf("Expected custom HTML, got %s", w.Body.String())
	}

	if w.Header().Get("Content-Type") != "text/html" {
		t.Errorf("Expected Content-Type to be text/html, got %s", w.Header().Get("Content-Type"))
	}
}

func TestCustomLoginRendererOnEmailError(t *testing.T) {
	customHTML := "<html><body>Custom Login Page with Error</body></html>"
	
	handler := New(
		LoginRenderer(func(data Attempt) (string, error) {
			if data.Error == "" {
				t.Errorf("Expected error message to be present")
			}
			return customHTML, nil
		}),
		EmailValidator(func(email string) error {
			return &ValidationError{Message: "Invalid email"}
		}),
	)

	form := url.Values{}
	form.Add("email", "invalid@example.com")
	
	req := httptest.NewRequest("POST", "/email", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusUnprocessableEntity {
		t.Errorf("Expected status code 422, got %d", w.Code)
	}

	if w.Body.String() != customHTML {
		t.Errorf("Expected custom HTML, got %s", w.Body.String())
	}
}

func TestFallbackToDefaultTemplates(t *testing.T) {
	// Test that default templates are used when no custom renderer is provided
	handler := New()

	req := httptest.NewRequest("GET", "/login", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", w.Code)
	}

	// Should contain default template content
	body := w.Body.String()
	if !strings.Contains(body, "Sign in to your account") {
		t.Errorf("Expected default login template content, got %s", body)
	}
}

type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

func TestCustomRenderersWithPrefix(t *testing.T) {
	customHTML := "<html><body>Custom Login with Prefix</body></html>"
	
	handler := New(
		Prefix("/auth/"),
		LoginRenderer(func(data Attempt) (string, error) {
			return customHTML, nil
		}),
	)

	req := httptest.NewRequest("GET", "/auth/login", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", w.Code)
	}

	if w.Body.String() != customHTML {
		t.Errorf("Expected custom HTML, got %s", w.Body.String())
	}
}

func TestCustomRendererError(t *testing.T) {
	handler := New(
		LoginRenderer(func(data Attempt) (string, error) {
			return "", &ValidationError{Message: "Renderer error"}
		}),
	)

	req := httptest.NewRequest("GET", "/login", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code 500, got %d", w.Code)
	}

	body := w.Body.String()
	if !strings.Contains(body, "Internal Server Error") {
		t.Errorf("Expected error response, got %s", body)
	}
}