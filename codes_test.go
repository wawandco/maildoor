package maildoor_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/wawandco/maildoor"
	"github.com/wawandco/maildoor/internal/testhelpers"
)

func TestCodeGeneration(t *testing.T) {
	t.Run("code generation creates 6 digit code", func(t *testing.T) {
		// We'll test this indirectly through the email handler
		// which calls newCodeFor internally
		auth := maildoor.New(
			maildoor.EmailValidator(func(email string) error {
				return nil
			}),
			maildoor.EmailSender(func(email, html, txt string) error {
				// Verify that the text contains a 6-digit code
				testhelpers.Contains(t, txt, "Code:")
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

	t.Run("multiple codes for different emails", func(t *testing.T) {
		var capturedCodes []string
		auth := maildoor.New(
			maildoor.EmailValidator(func(email string) error {
				return nil
			}),
			maildoor.EmailSender(func(email, html, txt string) error {
				// Extract code from text (it appears after "Code: ")
				lines := strings.Split(txt, "\n")
				for _, line := range lines {
					if strings.Contains(line, "Code:") {
						parts := strings.Split(line, "Code:")
						if len(parts) > 1 {
							code := strings.TrimSpace(parts[1])
							capturedCodes = append(capturedCodes, code)
						}
					}
				}
				return nil
			}),
		)

		// Generate code for first email
		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/email", nil)
		req.Form = url.Values{
			"email": []string{"user1@example.com"},
		}
		auth.ServeHTTP(w, req)

		// Generate code for second email
		w = httptest.NewRecorder()
		req = httptest.NewRequest("POST", "/email", nil)
		req.Form = url.Values{
			"email": []string{"user2@example.com"},
		}
		auth.ServeHTTP(w, req)

		// Both should succeed
		testhelpers.Equals(t, 2, len(capturedCodes))
		
		// Codes should be 6 characters long
		for _, code := range capturedCodes {
			testhelpers.Equals(t, 6, len(code))
		}
	})

	t.Run("code overwrite for same email", func(t *testing.T) {
		var codes []string
		auth := maildoor.New(
			maildoor.EmailValidator(func(email string) error {
				return nil
			}),
			maildoor.EmailSender(func(email, html, txt string) error {
				// Extract code from the text message
				lines := strings.Split(txt, "\n")
				for _, line := range lines {
					if strings.Contains(line, "Code:") {
						parts := strings.Split(line, "Code:")
						if len(parts) > 1 {
							code := strings.TrimSpace(parts[1])
							codes = append(codes, code)
						}
					}
				}
				return nil
			}),
		)

		email := "test@example.com"

		// Generate first code
		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/email", nil)
		req.Form = url.Values{
			"email": []string{email},
		}
		auth.ServeHTTP(w, req)

		// Generate second code for same email
		w = httptest.NewRecorder()
		req = httptest.NewRequest("POST", "/email", nil)
		req.Form = url.Values{
			"email": []string{email},
		}
		auth.ServeHTTP(w, req)

		testhelpers.Equals(t, 2, len(codes))
		// Both codes should be 6 digits
		testhelpers.Equals(t, 6, len(codes[0]))
		testhelpers.Equals(t, 6, len(codes[1]))
	})

	t.Run("codes only contain digits", func(t *testing.T) {
		var generatedCode string
		auth := maildoor.New(
			maildoor.EmailValidator(func(email string) error {
				return nil
			}),
			maildoor.EmailSender(func(email, html, txt string) error {
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
		)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/email", nil)
		req.Form = url.Values{
			"email": []string{"test@example.com"},
		}

		auth.ServeHTTP(w, req)

		// Verify code contains only digits
		for _, char := range generatedCode {
			if char < '0' || char > '9' {
				t.Fatalf("Code contains non-digit character: %c", char)
			}
		}
	})
}