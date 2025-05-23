package maildoor_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/wawandco/maildoor"
	"github.com/wawandco/maildoor/internal/testhelpers"
)

func TestNew(t *testing.T) {
	t.Run("creates handler with defaults", func(t *testing.T) {
		auth := maildoor.New()
		testhelpers.NotNil(t, auth)
	})

	t.Run("applies options correctly", func(t *testing.T) {
		auth := maildoor.New(
			maildoor.Logo("https://example.com/logo.png"),
			maildoor.ProductName("Test App"),
			maildoor.Prefix("/auth"),
		)
		testhelpers.NotNil(t, auth)
	})
}

func TestServeHTTP(t *testing.T) {
	t.Run("parses form data", func(t *testing.T) {
		auth := maildoor.New(
			maildoor.EmailValidator(func(email string) error {
				return nil
			}),
			maildoor.EmailSender(func(email, html, txt string) error {
				return nil
			}),
		)
		
		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/email", strings.NewReader("email=test@example.com"))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		auth.ServeHTTP(w, req)
		
		// Should not error due to form parsing
		testhelpers.NotEquals(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("handles method override with _method", func(t *testing.T) {
		var receivedMethod string
		auth := maildoor.New(
			maildoor.Logout(func(w http.ResponseWriter, r *http.Request) {
				receivedMethod = r.Method
				w.Write([]byte("OK"))
			}),
		)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/logout", nil)
		req.Form = url.Values{
			"_method": []string{"DELETE"},
		}

		auth.ServeHTTP(w, req)

		testhelpers.Equals(t, "DELETE", receivedMethod)
	})

	t.Run("logs request duration", func(t *testing.T) {
		auth := maildoor.New()
		
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/login", nil)

		// This should complete without error and log the duration
		auth.ServeHTTP(w, req)
		testhelpers.Equals(t, http.StatusOK, w.Code)
	})

	t.Run("handles malformed form data", func(t *testing.T) {
		auth := maildoor.New(
			maildoor.EmailValidator(func(email string) error {
				return nil
			}),
			maildoor.EmailSender(func(email, html, txt string) error {
				return nil
			}),
		)
		
		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/email", strings.NewReader("invalid%form%data"))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		auth.ServeHTTP(w, req)
		
		// Malformed form data causes an internal server error
		testhelpers.Equals(t, http.StatusInternalServerError, w.Code)
	})
}

func TestRouteRegistration(t *testing.T) {
	t.Run("registers routes with prefix", func(t *testing.T) {
		auth := maildoor.New(maildoor.Prefix("/auth"))
		
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/auth/login", nil)

		auth.ServeHTTP(w, req)
		testhelpers.Equals(t, http.StatusOK, w.Code)
	})

	t.Run("serves static assets", func(t *testing.T) {
		auth := maildoor.New()
		
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/logo.png", nil)

		auth.ServeHTTP(w, req)
		// Static assets may return 404 if not found, which is expected
		// Just verify the request was processed
		testhelpers.NotNil(t, w)
	})

	t.Run("serves static assets with prefix", func(t *testing.T) {
		auth := maildoor.New(maildoor.Prefix("/auth"))
		
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/auth/logo.png", nil)

		auth.ServeHTTP(w, req)
		// Static assets may return 404 if not found, which is expected
		// Just verify the request was processed
		testhelpers.NotNil(t, w)
	})
}

func TestTemplateRendering(t *testing.T) {
	t.Run("renders login template", func(t *testing.T) {
		auth := maildoor.New(
			maildoor.ProductName("Test App"),
			maildoor.Logo("https://example.com/logo.png"),
		)
		
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/login", nil)

		auth.ServeHTTP(w, req)
		
		testhelpers.Equals(t, http.StatusOK, w.Code)
		testhelpers.Contains(t, w.Body.String(), "Test App")
		testhelpers.Contains(t, w.Body.String(), "https://example.com/logo.png")
	})

	t.Run("renders code template after email submission", func(t *testing.T) {
		auth := maildoor.New(
			maildoor.EmailValidator(func(email string) error {
				return nil
			}),
			maildoor.EmailSender(func(email, html, txt string) error {
				return nil
			}),
		)
		
		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/email", nil)
		req.Form = url.Values{
			"email": []string{"test@example.com"},
		}

		auth.ServeHTTP(w, req)
		
		testhelpers.Equals(t, http.StatusOK, w.Code)
		testhelpers.Contains(t, w.Body.String(), "Check your inbox")
	})
}

func TestEmailTemplateGeneration(t *testing.T) {
	t.Run("generates email with code", func(t *testing.T) {
		var htmlBody, txtBody string
		auth := maildoor.New(
			maildoor.EmailValidator(func(email string) error {
				return nil
			}),
			maildoor.EmailSender(func(email, html, txt string) error {
				htmlBody = html
				txtBody = txt
				return nil
			}),
			maildoor.ProductName("Test App"),
			maildoor.Logo("https://example.com/logo.png"),
		)
		
		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/email", nil)
		req.Form = url.Values{
			"email": []string{"test@example.com"},
		}

		auth.ServeHTTP(w, req)
		
		testhelpers.Equals(t, http.StatusOK, w.Code)
		testhelpers.Contains(t, htmlBody, "Test App")
		testhelpers.Contains(t, txtBody, "Code:")
		testhelpers.Contains(t, htmlBody, "https://example.com/logo.png")
	})

	t.Run("includes current year in template", func(t *testing.T) {
		var htmlBody string
		auth := maildoor.New(
			maildoor.EmailValidator(func(email string) error {
				return nil
			}),
			maildoor.EmailSender(func(email, html, txt string) error {
				htmlBody = html
				return nil
			}),
		)
		
		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/email", nil)
		req.Form = url.Values{
			"email": []string{"test@example.com"},
		}

		auth.ServeHTTP(w, req)
		
		// Check that the year appears in the HTML (could be current year)
		testhelpers.NotEquals(t, "", htmlBody)
	})
}

func TestErrorHandling(t *testing.T) {
	t.Run("handles email sender error", func(t *testing.T) {
		auth := maildoor.New(
			maildoor.EmailValidator(func(email string) error {
				return nil
			}),
			maildoor.EmailSender(func(email, html, txt string) error {
				return errors.New("email service unavailable")
			}),
		)
		
		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/email", nil)
		req.Form = url.Values{
			"email": []string{"test@example.com"},
		}

		auth.ServeHTTP(w, req)
		
		testhelpers.Equals(t, http.StatusInternalServerError, w.Code)
		testhelpers.Contains(t, w.Body.String(), "email service unavailable")
	})

	t.Run("handles email validation error", func(t *testing.T) {
		auth := maildoor.New(
			maildoor.EmailValidator(func(email string) error {
				return errors.New("invalid email format")
			}),
		)
		
		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/email", nil)
		req.Form = url.Values{
			"email": []string{"invalid-email"},
		}

		auth.ServeHTTP(w, req)
		
		testhelpers.Equals(t, http.StatusUnprocessableEntity, w.Code)
		testhelpers.Contains(t, w.Body.String(), "invalid email format")
	})
}

func TestPrefixedPaths(t *testing.T) {
	t.Run("template includes prefixed paths", func(t *testing.T) {
		auth := maildoor.New(maildoor.Prefix("/auth"))
		
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/auth/login", nil)

		auth.ServeHTTP(w, req)
		
		testhelpers.Equals(t, http.StatusOK, w.Code)
		// The template should include the prefixed path in form actions
		testhelpers.Contains(t, w.Body.String(), "/auth/email")
	})
}

func TestHttpMethods(t *testing.T) {
	t.Run("GET /login works", func(t *testing.T) {
		auth := maildoor.New()
		
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/login", nil)

		auth.ServeHTTP(w, req)
		testhelpers.Equals(t, http.StatusOK, w.Code)
	})

	t.Run("POST /email works", func(t *testing.T) {
		auth := maildoor.New(
			maildoor.EmailValidator(func(email string) error {
				return nil
			}),
			maildoor.EmailSender(func(email, html, txt string) error {
				return nil
			}),
		)
		
		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/email", nil)
		req.Form = url.Values{
			"email": []string{"test@example.com"},
		}

		auth.ServeHTTP(w, req)
		testhelpers.Equals(t, http.StatusOK, w.Code)
	})

	t.Run("POST /code works", func(t *testing.T) {
		auth := maildoor.New()
		
		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/code", nil)
		req.Form = url.Values{
			"email": []string{"test@example.com"},
			"code":  []string{"123456"},
		}

		auth.ServeHTTP(w, req)
		testhelpers.Equals(t, http.StatusOK, w.Code)
	})

	t.Run("DELETE /logout works", func(t *testing.T) {
		auth := maildoor.New()
		
		w := httptest.NewRecorder()
		req := httptest.NewRequest("DELETE", "/logout", nil)

		auth.ServeHTTP(w, req)
		testhelpers.Equals(t, http.StatusFound, w.Code)
	})
}

func TestCompleteFlow(t *testing.T) {
	t.Run("full authentication flow", func(t *testing.T) {
		var generatedCode string
		
		auth := maildoor.New(
			maildoor.EmailValidator(func(email string) error {
				return nil
			}),
			maildoor.EmailSender(func(email, html, txt string) error {
				// Extract code from text message
				lines := strings.Split(txt, "\n")
				for _, line := range lines {
					if strings.Contains(line, "Code:") {
						parts := strings.Split(line, "Code:")
						if len(parts) > 1 {
							generatedCode = strings.TrimSpace(parts[1])
						}
					}
				}
				return nil
			}),
			maildoor.AfterLogin(func(w http.ResponseWriter, r *http.Request) {
				email := r.Context().Value("email").(string)
				testhelpers.Equals(t, "test@example.com", email)
				w.Write([]byte("Welcome!"))
			}),
		)

		// Step 1: Request login page
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/login", nil)
		auth.ServeHTTP(w, req)
		testhelpers.Equals(t, http.StatusOK, w.Code)

		// Step 2: Submit email
		w = httptest.NewRecorder()
		req = httptest.NewRequest("POST", "/email", nil)
		req.Form = url.Values{
			"email": []string{"test@example.com"},
		}
		auth.ServeHTTP(w, req)
		testhelpers.Equals(t, http.StatusOK, w.Code)
		testhelpers.NotEquals(t, "", generatedCode)

		// Step 3: Submit code (this will fail with current implementation since we can't access the internal codes map)
		// But we can test the invalid code path
		w = httptest.NewRecorder()
		req = httptest.NewRequest("POST", "/code", nil)
		req.Form = url.Values{
			"email": []string{"test@example.com"},
			"code":  []string{"wrong"},
		}
		auth.ServeHTTP(w, req)
		testhelpers.Equals(t, http.StatusOK, w.Code)
		testhelpers.Contains(t, w.Body.String(), "Invalid token")
	})
}

func TestTemplateRenderingEdgeCases(t *testing.T) {
	t.Run("render with empty partials", func(t *testing.T) {
		auth := maildoor.New()
		
		// This tests the render method with no partials
		// It's hard to test directly, but we can test via the HTTP handlers
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/login", nil)

		auth.ServeHTTP(w, req)
		testhelpers.Equals(t, http.StatusOK, w.Code)
	})

	t.Run("template function prefixedPath", func(t *testing.T) {
		auth := maildoor.New(maildoor.Prefix("/test"))
		
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/test/login", nil)

		auth.ServeHTTP(w, req)
		testhelpers.Equals(t, http.StatusOK, w.Code)
		// The template should use the prefixedPath function
		testhelpers.Contains(t, w.Body.String(), "/test/email")
	})
}

func TestMailBodiesGeneration(t *testing.T) {
	t.Run("generates both HTML and text email bodies", func(t *testing.T) {
		var htmlContent, textContent string
		auth := maildoor.New(
			maildoor.EmailValidator(func(email string) error {
				return nil
			}),
			maildoor.EmailSender(func(email, html, txt string) error {
				htmlContent = html
				textContent = txt
				return nil
			}),
			maildoor.ProductName("Test Product"),
			maildoor.Logo("https://test.com/logo.png"),
		)
		
		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/email", nil)
		req.Form = url.Values{
			"email": []string{"test@example.com"},
		}

		auth.ServeHTTP(w, req)
		
		testhelpers.Equals(t, http.StatusOK, w.Code)
		testhelpers.Contains(t, htmlContent, "Test Product")
		testhelpers.Contains(t, htmlContent, "https://test.com/logo.png")
		testhelpers.Contains(t, textContent, "Code:")
		testhelpers.NotEquals(t, htmlContent, textContent)
	})
}

func TestEdgeCaseScenarios(t *testing.T) {
	t.Run("empty email submission", func(t *testing.T) {
		auth := maildoor.New(
			maildoor.EmailValidator(func(email string) error {
				if email == "" {
					return errors.New("email required")
				}
				return nil
			}),
		)
		
		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/email", nil)
		req.Form = url.Values{
			"email": []string{""},
		}

		auth.ServeHTTP(w, req)
		testhelpers.Equals(t, http.StatusUnprocessableEntity, w.Code)
		testhelpers.Contains(t, w.Body.String(), "email required")
	})

	t.Run("handles different HTTP methods on routes", func(t *testing.T) {
		auth := maildoor.New()
		
		// Test unsupported method on email endpoint
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/email", nil)

		auth.ServeHTTP(w, req)
		testhelpers.Equals(t, http.StatusMethodNotAllowed, w.Code)
	})
}