package maildoor

import "net/http"

// Option is a function that can be passed to NewWithOptions to customize the
// behavior of the handler.
type Option func(*handler)

// UseProductName to be used in emails and login form.
func UseProductName(name string) Option {
	return func(h *handler) {
		h.product.Name = name
	}
}

// UseLogo for the login form and email.
func UseLogo(logoURL string) Option {
	return func(h *handler) {
		h.product.LogoURL = logoURL
	}
}

// UseFavicon for the login form and email.
func UseFavicon(faviconURL string) Option {
	return func(h *handler) {
		h.product.FaviconURL = faviconURL
	}
}

// UseLogger across the lifecycle of the handler.
func UseLogger(logger Logger) Option {
	return func(h *handler) {
		h.logger = logger
	}
}

// UsePrefix sets the prefix for the handler, this is
// useful for links and mounting the handler.
func UsePrefix(prefix string) Option {
	return func(h *handler) {
		h.prefix = prefix
	}
}

// UseBaseURL for links
func UseBaseURL(baseURL string) Option {
	return func(h *handler) {
		h.baseURL = baseURL
	}
}

// UseSender Specify the sender to be used by the handler.
func UseSender(fn func(message *Message) error) Option {
	return func(h *handler) {
		h.senderFn = fn
	}
}

// UseFinderFn sets the finder to be used.
func UseFinder(fn func(token string) (Emailable, error)) Option {
	return func(h *handler) {
		h.finderFn = fn
	}
}

// UseAfterLogin sets the function to be called after a successful login.
func UseAfterLogin(fn func(w http.ResponseWriter, r *http.Request, user Emailable) error) Option {
	return func(h *handler) {
		h.afterLoginFn = fn
	}
}

// UseLogout sets the function to be called after a successful logout.
func UseLogout(fn func(w http.ResponseWriter, r *http.Request) error) Option {
	return func(h *handler) {
		h.logoutFn = fn
	}
}

// UseTokenManager sets the token manager to be used.
func UseTokenManager(tokenManager TokenManager) Option {
	return func(h *handler) {
		h.tokenManager = tokenManager
	}
}
