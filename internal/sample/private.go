package sample

import (
	_ "embed"
	"net/http"
)

//go:embed private.html
var private []byte

// Private handler to show the private content to the user
func Private(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie("sample")
	if err != nil {
		http.Redirect(w, r, "/auth/login", http.StatusFound)
		return
	}

	w.Write(private)
}
