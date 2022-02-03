package maildoor

import "net/http"

// Options for the handler, these define the behavior of the
// handler while running, cannot be changed after initialized.
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

// Product options allow to customize the product name and logo
// as well as the favicon. These are used in the email that gets
// sent to the user and the login form.
type Product struct {
	Name       string
	LogoURL    string
	FaviconURL string
}
