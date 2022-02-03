package maildoor_test

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/wawandco/maildoor"
	"github.com/wawandco/maildoor/testhelpers"
)

func TestLogin(t *testing.T) {
	h, err := maildoor.New(maildoor.Options{
		CSRFTokenSecret: "secret",
	})

	testhelpers.NoError(t, err)

	t.Run("Content", func(tt *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/auth/login/", nil)
		w := httptest.NewRecorder()

		h.ServeHTTP(w, req)
		testhelpers.Equals(tt, http.StatusOK, w.Code)

		content := w.Body.String()

		testhelpers.Contains(tt, content, "Welcome Back 👋")
		testhelpers.Contains(tt, content, "/auth/send")
		testhelpers.Contains(tt, content, `<input type="hidden" name="CSRFToken" value="`)

		re := regexp.MustCompile("<input type=\"hidden\" name=\"CSRFToken\" value=\"(.*)\"")
		csrfToken := re.FindStringSubmatch(content)[1]
		testhelpers.NotEquals(t, len(csrfToken), 0)

		valid, err := maildoor.ValidateJWT(csrfToken, []byte("secret"))
		testhelpers.NoError(t, err)
		testhelpers.True(t, valid)
	})

}
