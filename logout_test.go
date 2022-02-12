package maildoor_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/wawandco/maildoor"
	"github.com/wawandco/maildoor/testhelpers"
)

func TestLogout(t *testing.T) {
	h, err := maildoor.New(maildoor.Options{
		CSRFTokenSecret: "secret",
		LogoutFn: func(w http.ResponseWriter, r *http.Request) error {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return nil
		},
	})

	testhelpers.NoError(t, err)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodDelete, "/auth/logout/", nil)

	h.ServeHTTP(w, r)
	testhelpers.Equals(t, http.StatusSeeOther, w.Code)
	testhelpers.Equals(t, w.Header().Get("Location"), "/")
}
