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
		auth := maildoor.New(maildoor.Prefix("/auth"))
		req := httptest.NewRequest("GET", "/auth/login", nil)
		w := httptest.NewRecorder()

		auth.ServeHTTP(w, req)

		testhelpers.Equals(t, http.StatusOK, w.Code)
		testhelpers.Contains(t, w.Body.String(), "Sign in to your account")
	})

	t.Run("using logo", func(t *testing.T) {
		auth := maildoor.New(
			maildoor.Prefix("/auth"),
			maildoor.Logo("https://my.logo/image.png"),
		)
		req := httptest.NewRequest("GET", "/auth/login", nil)
		w := httptest.NewRecorder()

		auth.ServeHTTP(w, req)

		testhelpers.Equals(t, http.StatusOK, w.Code)
		testhelpers.Contains(t, w.Body.String(), "Sign in to your account")
		testhelpers.Contains(t, w.Body.String(), "https://my.logo/image.png")
	})

	t.Run("using icon", func(t *testing.T) {
		auth := maildoor.New(
			maildoor.Prefix("/auth"),
			maildoor.Icon("https://my.icon/image.png"),
		)
		req := httptest.NewRequest("GET", "/auth/login", nil)
		w := httptest.NewRecorder()

		auth.ServeHTTP(w, req)

		testhelpers.Equals(t, http.StatusOK, w.Code)
		testhelpers.Contains(t, w.Body.String(), "Sign in to your account")
		testhelpers.Contains(t, w.Body.String(), "https://my.icon/image.png")
	})

	t.Run("using product name", func(t *testing.T) {
		auth := maildoor.New(
			maildoor.Prefix("/auth"),
			maildoor.ProductName("My App"),
		)
		req := httptest.NewRequest("GET", "/auth/login", nil)
		w := httptest.NewRecorder()

		auth.ServeHTTP(w, req)

		testhelpers.Equals(t, http.StatusOK, w.Code)
		testhelpers.Contains(t, w.Body.String(), "Sign in to your account")
		testhelpers.Contains(t, w.Body.String(), "My App")
	})

	t.Run("using all options", func(t *testing.T) {
		auth := maildoor.New(
			maildoor.Prefix("/auth"),
			maildoor.Logo("https://my.logo/image.png"),
			maildoor.Icon("https://my.icon/image.png"),
			maildoor.ProductName("My App"),
		)
		req := httptest.NewRequest("GET", "/auth/login", nil)
		w := httptest.NewRecorder()

		auth.ServeHTTP(w, req)

		testhelpers.Equals(t, http.StatusOK, w.Code)
		testhelpers.Contains(t, w.Body.String(), "Sign in to your account")
		testhelpers.Contains(t, w.Body.String(), "https://my.logo/image.png")
		testhelpers.Contains(t, w.Body.String(), "https://my.icon/image.png")
		testhelpers.Contains(t, w.Body.String(), "My App")
	})
}
