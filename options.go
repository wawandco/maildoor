package maildoor

import "net/http"

// Options for the handler constructor.
type Options struct {
	Prefix  string
	BaseURL string
	Product Product

	FinderFn     func(token string) (Emailable, error)
	SenderFn     func(message *Message) error
	AfterLoginFn func(w http.ResponseWriter, r *http.Request, user Emailable) error
	LogoutFn     func(w http.ResponseWriter, r *http.Request) error

	TokenManager TokenManager
}

type Product struct {
	Name       string
	LogoURL    string
	FaviconURL string
}
