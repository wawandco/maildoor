package internal

import (
	"html/template"
	"net/http"
)

// Login enpoint renders the login page to enter the user
// identifier.
func Login(w http.ResponseWriter, r *http.Request){
	page := struct{
		Logo string
	}{}


	tt, err := template.ParseFS(templates, "layout.html", "login.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = tt.Execute(w, page)
	if err != nil {
		return
	}
}
