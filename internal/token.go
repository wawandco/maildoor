package internal

import (
	"html/template"
	"net/http"
	"sync"
)

var tux sync.Mutex
var tokens = map[string]string{}

func Token(w http.ResponseWriter, r *http.Request){
	email := r.FormValue("email")

	tux.Lock()
	defer tux.Unlock()
	tokens[email] = "12345" // Token is generated here

	page := struct{
		Logo string
		Email string
	}{
		Logo: "",
		Email: email,
	}

	// Generate a token and store it in the server
	// Save the user email in the session
	tt, err := template.ParseFS(templates, "layout.html", "token.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = tt.Execute(w, page)
	if err != nil {
		return
	}
}
