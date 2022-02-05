package maildoor_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/wawandco/maildoor"
	"github.com/wawandco/maildoor/testhelpers"
)

func TestValidate(t *testing.T) {

	t.Run("Everything Valid", func(tt *testing.T) {
		var alcld = false
		h, err := maildoor.New(maildoor.Options{
			CSRFTokenSecret: "secret",
			FinderFn: func(token string) (maildoor.Emailable, error) {
				return testUser("amail@mail.com"), nil
			},

			AfterLoginFn: func(w http.ResponseWriter, r *http.Request, user maildoor.Emailable) error {
				alcld = true
				return nil
			},
		})

		testhelpers.NoError(t, err)

		token, err := maildoor.GenerateJWT(time.Second*10, []byte("not-so-secret-key"))
		testhelpers.NoError(t, err)

		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/auth/validate?token=%s&email=%v", token, "amil@amail.com"), nil)
		w := httptest.NewRecorder()

		h.ServeHTTP(w, req)
		testhelpers.NotEquals(tt, http.StatusInternalServerError, w.Code)
		testhelpers.True(t, alcld)
	})

	t.Run("expired token", func(tt *testing.T) {
		h, err := maildoor.New(maildoor.Options{
			CSRFTokenSecret: "secret",
			FinderFn: func(token string) (maildoor.Emailable, error) {
				return testUser("amail@mail.com"), nil
			},
		})

		testhelpers.NoError(t, err)

		token, err := maildoor.GenerateJWT(time.Second*-10, []byte("not-so-secret-key"))
		testhelpers.NoError(t, err)

		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/auth/validate?token=%s&email=%v", token, "amil@amail.com"), nil)
		w := httptest.NewRecorder()

		h.ServeHTTP(w, req)
		testhelpers.NotEquals(tt, http.StatusInternalServerError, w.Code)
		testhelpers.Equals(tt, http.StatusSeeOther, w.Code)
		testhelpers.Equals(tt, w.Header().Get("Location"), "http://127.0.0.1:8080/auth/login?error=E3")
	})

	t.Run("finder error", func(tt *testing.T) {
		h, err := maildoor.New(maildoor.Options{
			CSRFTokenSecret: "secret",
			FinderFn: func(token string) (maildoor.Emailable, error) {
				return nil, fmt.Errorf("error finding")
			},
		})

		testhelpers.NoError(t, err)

		token, err := maildoor.GenerateJWT(time.Second*10, []byte("not-so-secret-key"))
		testhelpers.NoError(t, err)

		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/auth/validate?token=%s&email=%v", token, "amil@amail.com"), nil)
		w := httptest.NewRecorder()

		h.ServeHTTP(w, req)
		testhelpers.NotEquals(tt, http.StatusInternalServerError, w.Code)
		testhelpers.Equals(tt, http.StatusSeeOther, w.Code)
		testhelpers.Equals(tt, w.Header().Get("Location"), "http://127.0.0.1:8080/auth/login?error=E1")
	})

	t.Run("nil user returned", func(tt *testing.T) {
		h, err := maildoor.New(maildoor.Options{
			CSRFTokenSecret: "secret",
			FinderFn: func(token string) (maildoor.Emailable, error) {
				return nil, nil
			},
		})

		testhelpers.NoError(t, err)

		token, err := maildoor.GenerateJWT(time.Second*10, []byte("not-so-secret-key"))
		testhelpers.NoError(t, err)

		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/auth/validate?token=%s&email=%v", token, "amil@amail.com"), nil)
		w := httptest.NewRecorder()

		h.ServeHTTP(w, req)
		testhelpers.NotEquals(tt, http.StatusInternalServerError, w.Code)
		testhelpers.Equals(tt, http.StatusSeeOther, w.Code)
		testhelpers.Equals(tt, w.Header().Get("Location"), "http://127.0.0.1:8080/auth/login?error=E7")
	})

	t.Run("nil user returned", func(tt *testing.T) {
		var alcld bool
		h, err := maildoor.New(maildoor.Options{
			CSRFTokenSecret: "secret",
			FinderFn: func(token string) (maildoor.Emailable, error) {
				return testUser("email@email.com"), nil
			},

			AfterLoginFn: func(w http.ResponseWriter, r *http.Request, user maildoor.Emailable) error {
				alcld = true
				return fmt.Errorf("error here")
			},
		})

		testhelpers.NoError(t, err)

		token, err := maildoor.GenerateJWT(time.Second*10, []byte("not-so-secret-key"))
		testhelpers.NoError(t, err)

		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/auth/validate?token=%s&email=%v", token, "amil@amail.com"), nil)
		w := httptest.NewRecorder()

		h.ServeHTTP(w, req)
		testhelpers.NotEquals(tt, http.StatusInternalServerError, w.Code)
		testhelpers.Equals(tt, http.StatusSeeOther, w.Code)
		testhelpers.True(tt, alcld)
		testhelpers.Equals(tt, w.Header().Get("Location"), "http://127.0.0.1:8080/auth/login?error=E7")
	})
}
