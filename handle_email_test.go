package maildoor_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/wawandco/maildoor"
	"github.com/wawandco/maildoor/internal/testhelpers"
)

func TestHandleEmail(t *testing.T) {
	// Test the handleEmail endpoint
	t.Run("basic test", func(t *testing.T) {
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
			"email": []string{"a@pagano.id"},
		}

		auth.ServeHTTP(w, req)

		testhelpers.Equals(t, http.StatusOK, w.Code)
		testhelpers.Contains(t, w.Body.String(), "Check your inbox")
	})

	t.Run("invalid email", func(t *testing.T) {
		auth := maildoor.New(
			maildoor.EmailValidator(func(email string) error {
				return errors.New("invalid email")
			}),
		)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/email", nil)
		req.Form = url.Values{
			"email": []string{"a@pagano.id"},
		}

		auth.ServeHTTP(w, req)
		testhelpers.Equals(t, http.StatusUnprocessableEntity, w.Code)
		testhelpers.Contains(t, w.Body.String(), "invalid email")
	})

	t.Run("error sending email", func(t *testing.T) {
		auth := maildoor.New(
			maildoor.EmailValidator(func(email string) error {
				return nil
			}),

			maildoor.EmailSender(func(email, html, txt string) error {
				return errors.New("error sending email")
			}),
		)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/email", nil)

		auth.ServeHTTP(w, req)
		testhelpers.Equals(t, http.StatusInternalServerError, w.Code)
		testhelpers.Contains(t, w.Body.String(), "error sending email")
	})

	t.Run("calls sending email", func(t *testing.T) {
		var textMessage string
		auth := maildoor.New(
			maildoor.EmailValidator(func(email string) error {
				return nil
			}),

			maildoor.EmailSender(func(email, html, txt string) error {
				textMessage = txt
				return nil
			}),
		)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/email", nil)

		auth.ServeHTTP(w, req)
		testhelpers.Equals(t, http.StatusOK, w.Code)
		testhelpers.Contains(t, textMessage, "Code:")
	})
}
