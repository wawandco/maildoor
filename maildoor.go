// package maildoor provides a passwordless authentication system which uses
// the email as the main authentication method.
package maildoor

import (
	"fmt"
	"net/http"
	"path"
)

type handler struct {
	prefix  string
	baseURL string
	product Product

	finderFn     func(token string) (Emailable, error)
	senderFn     func(message *Message) error
	afterLoginFn func(w http.ResponseWriter, r *http.Request, user Emailable) error
	logoutFn     func(w http.ResponseWriter, r *http.Request) error

	tokenManager TokenManager
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received:", r.URL.Path)

	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if r.URL.Path == "/login/" && r.Method == http.MethodGet {
		h.login(w, r)

		return
	}

	if r.URL.Path == "/send/" && r.Method == http.MethodPost {
		h.send(w, r)

		return
	}

	if r.URL.Path == "/validate/" && r.Method == http.MethodGet {
		h.validate(w, r)

		return
	}

	if r.URL.Path == "/logout/" && r.Method == http.MethodDelete {
		h.logout(w, r)

		return
	}

	http.NotFound(w, r)
}

func (h handler) sendPath() string {
	return h.baseURL + path.Join(h.prefix, "/send/")
}

func (h handler) loginPath() string {
	return h.baseURL + path.Join(h.prefix, "/login/")
}

func (h handler) validatePath() string {
	return h.baseURL + path.Join(h.prefix, "/validate/")
}

func New(o Options) *handler {
	h := &handler{}

	h.product = o.Product
	h.prefix = o.Prefix
	h.baseURL = o.BaseURL

	h.senderFn = o.SenderFn
	h.finderFn = o.FinderFn

	h.afterLoginFn = o.AfterLoginFn
	h.logoutFn = o.LogoutFn

	h.tokenManager = o.TokenManager

	return h
}
