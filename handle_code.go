package maildoor

import (
	"context"
	"html/template"
	"net/http"
	"strings"
)

// handleCode validates the input handleCode with the passed email.
func (m *maildoor) handleCode(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")

	// Find a combination of token and email in the server
	// call the afterlogin hook with the email
	// remove the token from the server
	token := strings.Join(r.Form["code[]"], "")
	valid := token == tokens[email]
	if !valid {
		// Generate a token and store it in the server
		// Save the user email in the session
		tt := template.New("layout.html").Funcs(m.helpers())
		tt, err := tt.ParseFS(templates, "layout.html", "handle_code.html")
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		err = tt.Execute(w, struct {
			Email string
			Error string
		}{
			Email: email,
			Error: "Invalid token",
		})

		if err != nil {
			return
		}

		return
	}

	delete(tokens, email)

	// Adding email to the context
	r = r.WithContext(context.WithValue(r.Context(), "email", email))
	m.afterLogin(w, r)
}
