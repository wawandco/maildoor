package internal

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"
)


func Validate(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	email:= r.FormValue("email")

	// Find a combination of token and email in the server
	// call the afterlogin hook with the email
	// remove the token from the server
	token := strings.Join(r.Form["code[]"], "")
	valid := token == tokens[email]
	if !valid {
		// Generate a token and store it in the server
		// Save the user email in the session
		tt, err := template.ParseFS(templates, "token_invalid.html")
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		err = tt.Execute(w, struct{
			Logo string
			Email string
		}{
			Email: email,
		})

		if err != nil {
			fmt.Println(err)
			return
		}

		return
	}

	delete(tokens, email)

	w.Write([]byte("Welcome"))
	// Do the after login hook
}
