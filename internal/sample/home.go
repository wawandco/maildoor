package sample

import (
	_ "embed"
	"net/http"
)

//go:embed home.html
var home []byte

// Simple home handler with link to the login page
func Home(w http.ResponseWriter, r *http.Request) {
	w.Write(home)
}
