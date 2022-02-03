package maildoor_test

import (
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
}
