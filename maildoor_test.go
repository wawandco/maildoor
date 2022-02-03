package maildoor

import (
	"testing"

	"github.com/wawandco/maildoor/testhelpers"
)

func TestNew(t *testing.T) {
	tcases := []struct {
		name   string
		opts   Options
		verify func(*testing.T, *handler)
	}{
		{
			name: "empty defaults",
			opts: Options{},
			verify: func(t *testing.T, h *handler) {
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
				Product: Product{
					Name:       "MyProduct",
					LogoURL:    "logoURL",
					FaviconURL: "faviconURL",
				},
			},
			verify: func(t *testing.T, h *handler) {
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

				tk, ok := h.tokenManager.(JWTTokenManager)
				testhelpers.True(t, ok)

				testhelpers.Equals(t, defaultTokenManager, tk)

			},
		},
	}

	for _, v := range tcases {
		t.Run(v.name, func(tt *testing.T) {
			h := New(v.opts)
			v.verify(tt, h)
		})
	}

}
