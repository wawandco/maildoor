package maildoor

import (
	"net/http"
)

// handler takes care of processing different actions against the maildoor
// server, such as login, send, validate, logout, most of these involve calling
// the corresponding functions provided by the host application.
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
