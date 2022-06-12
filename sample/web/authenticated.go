package web

import (
	"net/http"

	"github.com/wawandco/maildoor"
)

// authenticated is a middleware that checks if the user
// is logged in it redirects to the login page if not, its intended
// to be used with the routes that require authentication. It's
// also very simple in nature and wants to show how to keep some
// routes secure with maildoor.
func authenticated(h maildoor.CookieValuer, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		value, err := h.CookieValue(r)
		if err != nil {
			http.Redirect(w, r, "/auth/login/", http.StatusFound)
			return
		}

		u, err := finder(value)
		if u == nil || err != nil {
			http.Redirect(w, r, "/auth/login/", http.StatusFound)

			return
		}

		next(w, r)
	}
}
