package maildoor_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/wawandco/maildoor"
	"github.com/wawandco/maildoor/internal/testhelpers"
)

func TestHandleLogout(t *testing.T) {
	t.Run("default logout behavior", func(t *testing.T) {
		auth := maildoor.New()

		w := httptest.NewRecorder()
		req := httptest.NewRequest("DELETE", "/logout", nil)

		auth.ServeHTTP(w, req)

		testhelpers.Equals(t, http.StatusFound, w.Code)
		testhelpers.Equals(t, "/", w.Header().Get("Location"))
	})

	t.Run("custom logout handler", func(t *testing.T) {
		var logoutCalled bool
		auth := maildoor.New(
			maildoor.Logout(func(w http.ResponseWriter, r *http.Request) {
				logoutCalled = true
				w.Write([]byte("Logged out successfully"))
			}),
		)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("DELETE", "/logout", nil)

		auth.ServeHTTP(w, req)

		testhelpers.True(t, logoutCalled)
		testhelpers.Equals(t, http.StatusOK, w.Code)
		testhelpers.Contains(t, w.Body.String(), "Logged out successfully")
	})

	t.Run("logout with custom redirect", func(t *testing.T) {
		auth := maildoor.New(
			maildoor.Logout(func(w http.ResponseWriter, r *http.Request) {
				http.Redirect(w, r, "/login", http.StatusFound)
			}),
		)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("DELETE", "/logout", nil)

		auth.ServeHTTP(w, req)

		testhelpers.Equals(t, http.StatusFound, w.Code)
		testhelpers.Equals(t, "/login", w.Header().Get("Location"))
	})

	t.Run("logout with prefix", func(t *testing.T) {
		var logoutCalled bool
		auth := maildoor.New(
			maildoor.Prefix("/auth"),
			maildoor.Logout(func(w http.ResponseWriter, r *http.Request) {
				logoutCalled = true
				w.Write([]byte("Logged out"))
			}),
		)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("DELETE", "/auth/logout", nil)

		auth.ServeHTTP(w, req)

		testhelpers.True(t, logoutCalled)
		testhelpers.Equals(t, http.StatusOK, w.Code)
	})

	t.Run("logout via POST with _method override", func(t *testing.T) {
		var logoutCalled bool
		auth := maildoor.New(
			maildoor.Logout(func(w http.ResponseWriter, r *http.Request) {
				logoutCalled = true
				w.Write([]byte("Logged out via POST"))
			}),
		)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/logout", nil)
		req.Form = url.Values{
			"_method": []string{"DELETE"},
		}

		auth.ServeHTTP(w, req)

		testhelpers.True(t, logoutCalled)
		testhelpers.Equals(t, http.StatusOK, w.Code)
		testhelpers.Contains(t, w.Body.String(), "Logged out via POST")
	})

	t.Run("logout receives correct request", func(t *testing.T) {
		var receivedMethod string
		var receivedPath string
		auth := maildoor.New(
			maildoor.Logout(func(w http.ResponseWriter, r *http.Request) {
				receivedMethod = r.Method
				receivedPath = r.URL.Path
				w.Write([]byte("OK"))
			}),
		)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("DELETE", "/logout", nil)

		auth.ServeHTTP(w, req)

		testhelpers.Equals(t, "DELETE", receivedMethod)
		testhelpers.Equals(t, "/logout", receivedPath)
	})
}