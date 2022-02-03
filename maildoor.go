// package maildoor provides a passwordless authentication system which uses
// the email as the main authentication method.
package maildoor

import (
	"errors"
)

var (
	defaultPrefix       = "/auth"
	defaultBaseURL      = "http://127.0.0.1:8080"
	defaultTokenManager = JWTTokenManager("not-so-secret-key")

	defaultProduct = Product{
		Name:       "maildoor",
		LogoURL:    "https://maildoor.com/assets/images/logo.png",
		FaviconURL: "https://maildoor.com/assets/images/favicon.ico",
	}

	defaultSender = func(message *Message) error {
		return errors.New("did not send message")
	}

	defaultFinder = func(token string) (Emailable, error) {
		return nil, errors.New("did not find user")
	}
)

// New maildoor handler with the given options, all of the options have defaults,
// if not specified this method pulls the default value for them.
func New(o Options) *handler {
	h := &handler{
		product: defaultProduct,
		prefix:  defaultPrefix,
		baseURL: defaultBaseURL,

		senderFn:     defaultSender,
		finderFn:     defaultFinder,
		tokenManager: defaultTokenManager,
	}

	if o.Product != (Product{}) {
		h.product = o.Product
	}

	if o.Prefix != "" {
		h.prefix = o.Prefix
	}

	if o.BaseURL != "" {
		h.baseURL = o.BaseURL
	}

	if o.SenderFn != nil {
		h.senderFn = o.SenderFn
	}

	if o.FinderFn != nil {
		h.finderFn = o.FinderFn
	}

	if o.AfterLoginFn != nil {
		h.afterLoginFn = o.AfterLoginFn
	}

	if o.LogoutFn != nil {
		h.logoutFn = o.LogoutFn
	}

	if o.TokenManager != nil {
		h.tokenManager = o.TokenManager
	}

	return h
}
