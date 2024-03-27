package main

import (
	"log/slog"
	"net/http"

	"github.com/wawandco/maildoor/internal/sample"
)

func main() {
	r := http.NewServeMux()

	// Auth handlers
	r.Handle("/auth/", sample.Auth)

	// Application handlers
	r.HandleFunc("/private", secure(sample.Private))
	r.HandleFunc("/{$}", sample.Home)

	slog.Info("Server running on :3000")
	http.ListenAndServe(":3000", r)
}

// secure middleware checks if the user is authenticated
// if not, it redirects to the login page
func secure(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := r.Cookie("sample")
		if err != nil {
			http.Redirect(w, r, "/auth/login", http.StatusFound)
			return
		}

		next.ServeHTTP(w, r)
	}
}
