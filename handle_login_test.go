package maildoor_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/wawandco/maildoor"
	"github.com/wawandco/maildoor/internal/testhelpers"
)

func TestHandleLogin(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		auth := maildoor.New()
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/login", nil)

		auth.ServeHTTP(w, req)

		testhelpers.Equals(t, http.StatusOK, w.Code)
		testhelpers.Contains(t, w.Body.String(), "Sign in to your account")
	})

	t.Run("not extra path", func(t *testing.T) {
		auth := maildoor.New()
		req := httptest.NewRequest("GET", "/login/other-thing", nil)
		w := httptest.NewRecorder()

		auth.ServeHTTP(w, req)
		testhelpers.Equals(t, http.StatusNotFound, w.Code)
	})

	t.Run("using prefix", func(t *testing.T) {
		auth := maildoor.New(maildoor.UsePrefix("/auth"))
		req := httptest.NewRequest("GET", "/auth/login", nil)
		w := httptest.NewRecorder()

		auth.ServeHTTP(w, req)

		testhelpers.Equals(t, http.StatusOK, w.Code)
		testhelpers.Contains(t, w.Body.String(), "Sign in to your account")
	})

	t.Run("using logo", func(t *testing.T) {
		auth := maildoor.New(
			maildoor.UsePrefix("/auth"),
			maildoor.WithLogo("https://my.logo/image.png"),
		)
		req := httptest.NewRequest("GET", "/auth/login", nil)
		w := httptest.NewRecorder()

		auth.ServeHTTP(w, req)

		testhelpers.Equals(t, http.StatusOK, w.Code)
		testhelpers.Contains(t, w.Body.String(), "Sign in to your account")
		testhelpers.Contains(t, w.Body.String(), "https://my.logo/image.png")
	})
}
