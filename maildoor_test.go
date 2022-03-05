package maildoor

import (
	"testing"

	"github.com/wawandco/maildoor/internal/testhelpers"
)

func TestNew(t *testing.T) {
	tcases := []struct {
		name   string
		opts   Options
		verify func(*testing.T, *handler, error)
	}{
		{
			name: "totally empty it errors",
			opts: Options{},
			verify: func(t *testing.T, h *handler, err error) {
				testhelpers.Error(t, err)
			},
		},
		{
			name: "empty defaults",
			opts: Options{
				CSRFTokenSecret: "secret",
			},
			verify: func(t *testing.T, h *handler, err error) {
				testhelpers.NoError(t, err)
				testhelpers.NotEquals(t, "", h.product.Name)
				testhelpers.NotEquals(t, "", h.prefix)
				testhelpers.NotEquals(t, "", h.baseURL)
				testhelpers.NotNil(t, h.tokenManager)
				testhelpers.NotNil(t, h.afterLoginFn)
				testhelpers.NotNil(t, h.senderFn)
				testhelpers.NotNil(t, h.finderFn)
				testhelpers.NotNil(t, h.logoutFn)
			},
		},

		{
			name: "some empty defaults",
			opts: Options{
				CSRFTokenSecret: "secret",
				Product: Product{
					Name:       "MyProduct",
					LogoURL:    "logoURL",
					FaviconURL: "faviconURL",
				},
			},
			verify: func(t *testing.T, h *handler, err error) {
				testhelpers.NoError(t, err)

				testhelpers.Equals(t, "MyProduct", h.product.Name)
				testhelpers.Equals(t, "logoURL", h.product.LogoURL)
				testhelpers.Equals(t, "faviconURL", h.product.FaviconURL)

				testhelpers.NotEquals(t, "", h.prefix)
				testhelpers.NotEquals(t, "", h.baseURL)
				testhelpers.NotNil(t, h.tokenManager)
				testhelpers.NotNil(t, h.afterLoginFn)
				testhelpers.NotNil(t, h.senderFn)
				testhelpers.NotNil(t, h.finderFn)
				testhelpers.NotNil(t, h.logoutFn)

				tk, ok := h.tokenManager.(DefaultTokenManager)
				testhelpers.True(t, ok)

				testhelpers.Equals(t, defaultTokenManager, tk)

			},
		},
		{
			name: "product images",
			opts: Options{
				CSRFTokenSecret: "secret",
			},
			verify: func(t *testing.T, h *handler, err error) {
				testhelpers.NoError(t, err)
				testhelpers.NotEquals(t, "", h.product.Name)
				testhelpers.NotEquals(t, "", h.prefix)
				testhelpers.NotEquals(t, "", h.baseURL)
				testhelpers.NotNil(t, h.tokenManager)
				testhelpers.NotNil(t, h.afterLoginFn)
				testhelpers.NotNil(t, h.senderFn)
				testhelpers.NotNil(t, h.finderFn)
				testhelpers.NotNil(t, h.logoutFn)
			},
		},
	}

	for _, v := range tcases {
		t.Run(v.name, func(tt *testing.T) {
			h, err := New(v.opts)
			v.verify(tt, h, err)
		})
	}

}
