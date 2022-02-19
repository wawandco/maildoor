package web

import (
	"net/http"
)

// authenticated is a middleware that checks if the user
// is logged in it redirects to the login page if not, its intended
// to be used with the routes that require authentication. It's
// also very simple in nature and wants to show how to keep some
// routes secure with maildoor.
func authenticated(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("sample")
		if err != nil {
			http.Redirect(w, r, "/auth/login/", http.StatusFound)
			return
		}

		u, err := finder(c.Value)
		if u == nil || err != nil {
			http.Redirect(w, r, "/auth/login/", http.StatusFound)

			return
		}

		next(w, r)
	}
}
