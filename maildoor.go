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
	defaultProduct      = productConfig{
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

func NewWithOptions(csrfToken string, options ...Option) (*handler, error) {
	h := &handler{
		product: defaultProduct,
		prefix:  defaultPrefix,
		baseURL: defaultBaseURL,

		senderFn:     defaultSender,
		finderFn:     defaultFinder,
		tokenManager: defaultTokenManager,
		logger:       defaultLogger,

		assetsServer:    http.FileServer(http.FS(assets)),
		csrfTokenSecret: csrfToken,

		logoutFn:     defaultLogout,
		afterLoginFn: defaultAfterLogin,
	}

	if csrfToken == "" {
		return nil, errors.New("CSRF token is empty")
	}

	// Apply each of the options passed for the Maildoor instance.
	for _, option := range options {
		option(h)
	}

	return h, nil
}
