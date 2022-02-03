package maildoor_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/wawandco/maildoor"
	"github.com/wawandco/maildoor/testhelpers"
)

func TestHandlerLogin(t *testing.T) {
	h := maildoor.New(maildoor.Options{})

	req := httptest.NewRequest(http.MethodGet, "/auth/login/", nil)
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)
	testhelpers.Equals(t, http.StatusOK, w.Code)
}
