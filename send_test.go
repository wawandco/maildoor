package maildoor_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/wawandco/maildoor"
	"github.com/wawandco/maildoor/testhelpers"
)

func TestSend(t *testing.T) {

	t.Run("Invalid CSRF", func(tt *testing.T) {
		h, err := maildoor.New(maildoor.Options{
			CSRFTokenSecret: "secret",
		})

		testhelpers.NoError(t, err)

		data := url.Values{
			"CSRFToken": {"invalid"},
		}

		req := httptest.NewRequest(http.MethodPost, "/auth/send/", strings.NewReader(data.Encode()))
		w := httptest.NewRecorder()

		h.ServeHTTP(w, req)
		testhelpers.Equals(tt, http.StatusSeeOther, w.Code)
		testhelpers.Equals(tt, w.HeaderMap.Get("Location"), "http://127.0.0.1:8080/auth/login?error=E4")
	})

	t.Run("Valid CSRF", func(tt *testing.T) {
		h, err := maildoor.New(maildoor.Options{
			CSRFTokenSecret: "secret",
			FinderFn: func(token string) (maildoor.Emailable, error) {
				return nil, nil
			},
		})

		testhelpers.NoError(t, err)

		token, err := maildoor.GenerateJWT(time.Second*10, []byte("secret"))
		testhelpers.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/auth/send/", nil)
		req.Form = url.Values{
			"CSRFToken": []string{token},
		}

		w := httptest.NewRecorder()

		h.ServeHTTP(w, req)
		testhelpers.Equals(tt, http.StatusOK, w.Code)
	})

	t.Run("Valid Error Finding", func(tt *testing.T) {
		h, err := maildoor.New(maildoor.Options{
			CSRFTokenSecret: "secret",
			FinderFn: func(token string) (maildoor.Emailable, error) {
				return nil, fmt.Errorf("error finding the user")
			},
		})

		testhelpers.NoError(t, err)

		token, err := maildoor.GenerateJWT(time.Second*10, []byte("secret"))
		testhelpers.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/auth/send/", nil)
		req.Form = url.Values{
			"CSRFToken": []string{token},
			"email":     []string{"test@test.com"},
		}

		w := httptest.NewRecorder()

		h.ServeHTTP(w, req)
		testhelpers.Equals(tt, http.StatusSeeOther, w.Code)
		testhelpers.Equals(tt, w.Header().Get("Location"), "http://127.0.0.1:8080/auth/login?error=E1")
	})

	t.Run("User Not found renders ok", func(tt *testing.T) {
		h, err := maildoor.New(maildoor.Options{
			CSRFTokenSecret: "secret",
			FinderFn: func(token string) (maildoor.Emailable, error) {
				return nil, nil
			},
		})

		testhelpers.NoError(t, err)

		token, err := maildoor.GenerateJWT(time.Second*10, []byte("secret"))
		testhelpers.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/auth/send/", nil)
		req.Form = url.Values{
			"CSRFToken": []string{token},
			"email":     []string{"test@test.com"},
		}

		w := httptest.NewRecorder()

		h.ServeHTTP(w, req)
		testhelpers.Equals(tt, http.StatusOK, w.Code)
	})

	t.Run("User found", func(tt *testing.T) {
		var m maildoor.Message
		h, err := maildoor.New(maildoor.Options{
			CSRFTokenSecret: "secret",
			FinderFn: func(token string) (maildoor.Emailable, error) {
				return testUser("mailo@wawand.co"), nil
			},

			SenderFn: func(message *maildoor.Message) error {
				m = *message
				return nil
			},
		})

		testhelpers.NoError(t, err)

		token, err := maildoor.GenerateJWT(time.Second*10, []byte("secret"))
		testhelpers.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/auth/send/", nil)
		req.Form = url.Values{
			"CSRFToken": []string{token},
			"email":     []string{"mailo@wawand.com"},
		}

		w := httptest.NewRecorder()

		h.ServeHTTP(w, req)
		testhelpers.Equals(tt, http.StatusOK, w.Code)
		testhelpers.Equals(tt, "mailo@wawand.co", m.To)
		testhelpers.Contains(tt, string(m.Bodies[0].Content), "http://127.0.0.1:8080/auth/validate")
	})

	t.Run("User found error sending", func(tt *testing.T) {
		h, err := maildoor.New(maildoor.Options{
			CSRFTokenSecret: "secret",
			FinderFn: func(token string) (maildoor.Emailable, error) {
				return testUser("mailo@wawand.co"), nil
			},

			SenderFn: func(message *maildoor.Message) error {
				return fmt.Errorf("error sending")
			},
		})

		testhelpers.NoError(t, err)

		token, err := maildoor.GenerateJWT(time.Second*10, []byte("secret"))
		testhelpers.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/auth/send/", nil)
		req.Form = url.Values{
			"CSRFToken": []string{token},
			"email":     []string{"mailo@wawand.com"},
		}

		w := httptest.NewRecorder()

		h.ServeHTTP(w, req)
		testhelpers.Equals(tt, http.StatusSeeOther, w.Code)
		testhelpers.Equals(tt, w.Header().Get("Location"), "http://127.0.0.1:8080/auth/login?error=E5")
	})

	t.Run("User found error generating token", func(tt *testing.T) {
		h, err := maildoor.New(maildoor.Options{
			CSRFTokenSecret: "secret",
			FinderFn: func(token string) (maildoor.Emailable, error) {
				return testUser("mailo@wawand.co"), nil
			},

			SenderFn: func(message *maildoor.Message) error {
				return fmt.Errorf("error sending")
			},

			TokenManager: errTokenManager("error generating token"),
		})

		testhelpers.NoError(t, err)

		token, err := maildoor.GenerateJWT(time.Second*10, []byte("secret"))
		testhelpers.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/auth/send/", nil)
		req.Form = url.Values{
			"CSRFToken": []string{token},
			"email":     []string{"mailo@wawand.com"},
		}

		w := httptest.NewRecorder()

		h.ServeHTTP(w, req)
		testhelpers.Equals(tt, http.StatusSeeOther, w.Code)
		testhelpers.Equals(tt, w.Header().Get("Location"), "http://127.0.0.1:8080/auth/login?error=E6")
	})

}
