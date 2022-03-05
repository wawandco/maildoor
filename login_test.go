package maildoor_test

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/wawandco/maildoor"
	"github.com/wawandco/maildoor/internal/testhelpers"
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

	t.Run("Errors", func(tt *testing.T) {

		tcases := []struct {
			name string
			url  string
			val  func(t *testing.T, content string)
		}{
			{
				name: "E1",
				url:  "/auth/login/?error=E1",
				val: func(t *testing.T, content string) {
					testhelpers.Contains(tt, content, `😥  something happened while trying`)
				},
			},

			{
				name: "E2",
				url:  "/auth/login/?error=E2",
				val: func(t *testing.T, content string) {
					testhelpers.Contains(tt, content, `sorry, the specified token has already expired`)
				},
			},

			{
				name: "E3",
				url:  "/auth/login/?error=E3",
				val: func(t *testing.T, content string) {
					testhelpers.Contains(tt, content, `The token you have entered is invalid.`)
				},
			},

			{
				name: "No Error",
				url:  "/auth/login/",
				val: func(t *testing.T, content string) {
					testhelpers.NotContains(tt, content, `<div class="rounded bg-red-100 border border-red-300 py-3 px-8 text-red-500 mb-4" role="alert">`)
				},
			},

			{
				name: "Invalid Error",
				url:  "/auth/login/?err=SOMETHING",
				val: func(t *testing.T, content string) {
					testhelpers.NotContains(tt, content, `<div class="rounded bg-red-100 border border-red-300 py-3 px-8 text-red-500 mb-4" role="alert">`)
				},
			},
		}

		for _, tcase := range tcases {
			tt.Run(tcase.name, func(tx *testing.T) {
				req := httptest.NewRequest(http.MethodGet, tcase.url, nil)
				w := httptest.NewRecorder()
				h.ServeHTTP(w, req)

				testhelpers.Equals(tx, http.StatusOK, w.Code)
				tcase.val(tx, w.Body.String())
			})
		}

	})

}
