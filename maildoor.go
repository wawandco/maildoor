// Package maildoor provides a passwordless authentication system which uses
// the email as the main authentication method. It allows applications to configure
// its specific behavior and takes care of the email authentication behavior. It provides
// default implementations for token generation and validation as well as email sending
// through SMTP.
package maildoor

import (
	"embed"
	"errors"
	"net/http"
	"time"
)

var (
	// csrfDuration defines the lifetime of generated CSRF token (JWT).
	csrfDuration = 3 * time.Minute

	defaultPrefix       = "/auth"
	defaultBaseURL      = "http://127.0.0.1:8080"
	defaultTokenManager = DefaultTokenManager("not-so-secret-key")
	defaultProduct      = Product{
		Name:       "maildoor",
		LogoURL:    "https://github.com/wawandco/maildoor/raw/main/images/maildoor_logo.png",
		FaviconURL: "https://github.com/wawandco/maildoor/raw/main/images/favicon.png",
	}

	defaultSender = func(message *Message) error {
		return errors.New("did not send message")
	}

	defaultFinder = func(token string) (Emailable, error) {
		return nil, errors.New("did not find user")
	}

	//go:embed assets
	assets embed.FS
)

// New maildoor handler with the given options, all of the options have defaults,
// if not specified this method pulls the default value for them.
func New(o Options) (*handler, error) {
	h := &handler{
		product: defaultProduct,
		prefix:  defaultPrefix,
		baseURL: defaultBaseURL,

		senderFn:     defaultSender,
		finderFn:     defaultFinder,
		tokenManager: defaultTokenManager,
		logger:       defaultLogger,

		assetsServer: http.FileServer(http.FS(assets)),
	}

	h.product.LogoURL = h.logoPath()
	h.product.FaviconURL = h.faviconPath()

	if o.CSRFTokenSecret == "" {
		return nil, errors.New("CSRF Token secret is required")
	}

	h.csrfTokenSecret = o.CSRFTokenSecret

	if o.Product != (Product{}) {
		if o.Product.LogoURL == "" {
			o.Product.LogoURL = h.logoPath()
		}

		if o.Product.FaviconURL == "" {
			o.Product.FaviconURL = h.faviconPath()
		}

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

	if o.Logger != nil {
		h.logger = o.Logger
	}

	return h, nil
}
