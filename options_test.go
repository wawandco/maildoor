package maildoor_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/wawandco/maildoor"
	"github.com/wawandco/maildoor/internal/testhelpers"
)

func TestOptions(t *testing.T) {
	t.Run("Logo option", func(t *testing.T) {
		logoURL := "https://example.com/custom-logo.png"
		auth := maildoor.New(maildoor.Logo(logoURL))

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/login", nil)

		auth.ServeHTTP(w, req)

		testhelpers.Equals(t, http.StatusOK, w.Code)
		testhelpers.Contains(t, w.Body.String(), logoURL)
	})

	t.Run("ProductName option", func(t *testing.T) {
		productName := "My Custom App"
		auth := maildoor.New(maildoor.ProductName(productName))

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/login", nil)

		auth.ServeHTTP(w, req)

		testhelpers.Equals(t, http.StatusOK, w.Code)
		testhelpers.Contains(t, w.Body.String(), productName)
	})

	t.Run("Icon option", func(t *testing.T) {
		iconURL := "https://example.com/custom-icon.png"
		auth := maildoor.New(maildoor.Icon(iconURL))

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/login", nil)

		auth.ServeHTTP(w, req)

		testhelpers.Equals(t, http.StatusOK, w.Code)
		testhelpers.Contains(t, w.Body.String(), iconURL)
	})

	t.Run("Prefix option", func(t *testing.T) {
		prefix := "/custom-auth"
		auth := maildoor.New(maildoor.Prefix(prefix))

		// Test that the prefixed route works
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/custom-auth/login", nil)

		auth.ServeHTTP(w, req)

		testhelpers.Equals(t, http.StatusOK, w.Code)
	})

	t.Run("Prefix option - non-prefixed route fails", func(t *testing.T) {
		prefix := "/custom-auth"
		auth := maildoor.New(maildoor.Prefix(prefix))

		// Test that non-prefixed route fails
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/login", nil)

		auth.ServeHTTP(w, req)

		testhelpers.Equals(t, http.StatusNotFound, w.Code)
	})

	t.Run("AfterLogin option", func(t *testing.T) {
		auth := maildoor.New(
			maildoor.AfterLogin(func(w http.ResponseWriter, r *http.Request) {
				if email := r.Context().Value("email"); email != nil {
					_ = email.(string)
				}
				w.Write([]byte("Custom after login"))
			}),
		)

		// Test with invalid code to trigger the flow
		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/code", nil)
		req.Form = map[string][]string{
			"email": {"test@example.com"},
			"code":  {"invalid"},
		}

		auth.ServeHTTP(w, req)

		// This will show invalid token, but we can't easily test the success path
		testhelpers.Contains(t, w.Body.String(), "Invalid token")
	})

	t.Run("Logout option", func(t *testing.T) {
		var logoutCalled bool
		auth := maildoor.New(
			maildoor.Logout(func(w http.ResponseWriter, r *http.Request) {
				logoutCalled = true
				w.Write([]byte("Custom logout"))
			}),
		)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("DELETE", "/logout", nil)

		auth.ServeHTTP(w, req)

		testhelpers.True(t, logoutCalled)
		testhelpers.Contains(t, w.Body.String(), "Custom logout")
	})

	t.Run("EmailValidator option", func(t *testing.T) {
		var validatorCalled bool
		var receivedEmail string

		auth := maildoor.New(
			maildoor.EmailValidator(func(email string) error {
				validatorCalled = true
				receivedEmail = email
				return nil
			}),
			maildoor.EmailSender(func(email, html, txt string) error {
				return nil
			}),
		)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/email", nil)
		req.Form = map[string][]string{
			"email": {"test@example.com"},
		}

		auth.ServeHTTP(w, req)

		testhelpers.True(t, validatorCalled)
		testhelpers.Equals(t, "test@example.com", receivedEmail)
		testhelpers.Equals(t, http.StatusOK, w.Code)
	})

	t.Run("EmailSender option", func(t *testing.T) {
		var senderCalled bool
		var receivedEmail string
		var receivedHTML string
		var receivedTxt string

		auth := maildoor.New(
			maildoor.EmailValidator(func(email string) error {
				return nil
			}),
			maildoor.EmailSender(func(email, html, txt string) error {
				senderCalled = true
				receivedEmail = email
				receivedHTML = html
				receivedTxt = txt
				return nil
			}),
		)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/email", nil)
		req.Form = map[string][]string{
			"email": {"test@example.com"},
		}

		auth.ServeHTTP(w, req)

		testhelpers.True(t, senderCalled)
		testhelpers.Equals(t, "test@example.com", receivedEmail)
		testhelpers.NotEquals(t, "", receivedHTML)
		testhelpers.NotEquals(t, "", receivedTxt)
		testhelpers.Contains(t, receivedTxt, "Code:")
	})

	t.Run("multiple options combined", func(t *testing.T) {
		var validatorCalled bool
		var senderCalled bool

		auth := maildoor.New(
			maildoor.Logo("https://example.com/logo.png"),
			maildoor.ProductName("Multi-Option App"),
			maildoor.Icon("https://example.com/icon.png"),
			maildoor.Prefix("/auth"),
			maildoor.EmailValidator(func(email string) error {
				validatorCalled = true
				return nil
			}),
			maildoor.EmailSender(func(email, html, txt string) error {
				senderCalled = true
				return nil
			}),
		)

		// Test login page with all options
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/auth/login", nil)

		auth.ServeHTTP(w, req)

		testhelpers.Equals(t, http.StatusOK, w.Code)
		testhelpers.Contains(t, w.Body.String(), "https://example.com/logo.png")
		testhelpers.Contains(t, w.Body.String(), "Multi-Option App")
		testhelpers.Contains(t, w.Body.String(), "https://example.com/icon.png")

		// Test email submission
		w = httptest.NewRecorder()
		req = httptest.NewRequest("POST", "/auth/email", nil)
		req.Form = map[string][]string{
			"email": {"test@example.com"},
		}

		auth.ServeHTTP(w, req)

		testhelpers.True(t, validatorCalled)
		testhelpers.True(t, senderCalled)
		testhelpers.Equals(t, http.StatusOK, w.Code)
	})

	t.Run("default values when no options provided", func(t *testing.T) {
		auth := maildoor.New()

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/login", nil)

		auth.ServeHTTP(w, req)

		testhelpers.Equals(t, http.StatusOK, w.Code)
		// Should contain default product name
		testhelpers.Contains(t, w.Body.String(), "Maildoor")
		// Should contain default logo URL
		testhelpers.Contains(t, w.Body.String(), "maildoor_logo.png")
	})

	t.Run("empty prefix option", func(t *testing.T) {
		auth := maildoor.New(maildoor.Prefix(""))

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/login", nil)

		auth.ServeHTTP(w, req)

		testhelpers.Equals(t, http.StatusOK, w.Code)
	})

	t.Run("prefix with trailing slash", func(t *testing.T) {
		auth := maildoor.New(maildoor.Prefix("/auth/"))

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/auth/login", nil)

		auth.ServeHTTP(w, req)

		testhelpers.Equals(t, http.StatusOK, w.Code)
	})

	t.Run("very long product name", func(t *testing.T) {
		longName := "This is a very long product name that should still work correctly in the templates and not break anything"
		auth := maildoor.New(maildoor.ProductName(longName))

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/login", nil)

		auth.ServeHTTP(w, req)

		testhelpers.Equals(t, http.StatusOK, w.Code)
		testhelpers.Contains(t, w.Body.String(), longName)
	})

	t.Run("special characters in product name", func(t *testing.T) {
		specialName := "My App & Co. <Test>"
		auth := maildoor.New(maildoor.ProductName(specialName))

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/login", nil)

		auth.ServeHTTP(w, req)

		testhelpers.Equals(t, http.StatusOK, w.Code)
		// HTML should be properly escaped
		testhelpers.Contains(t, w.Body.String(), "My App &amp; Co.")
	})
}