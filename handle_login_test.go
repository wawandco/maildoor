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
		if w.Code != 200 {
			t.Errorf("Expected status code 200, got %d", w.Code)
		}

		testhelpers.Contains(t, w.Body.String(), "Welcome")
	})

	t.Run("not extra path", func(t *testing.T) {
		auth := maildoor.New()
		req := httptest.NewRequest("GET", "/login/other-thing", nil)
		w := httptest.NewRecorder()

		auth.ServeHTTP(w, req)
		if w.Code != http.StatusNotFound {
			t.Errorf("Expected status code 404, got %d", w.Code)
		}
	})

	t.Run("using prefix", func(t *testing.T) {
		auth := maildoor.New(maildoor.UsePrefix("/auth"))
		req := httptest.NewRequest("GET", "/auth/login", nil)
		w := httptest.NewRecorder()

		auth.ServeHTTP(w, req)
		if w.Code != 200 {
			t.Errorf("Expected status code 200, got %d", w.Code)
		}

		testhelpers.Contains(t, w.Body.String(), "Welcome")
	})
}
