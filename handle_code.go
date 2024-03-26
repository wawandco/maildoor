package maildoor

import (
	"context"
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
	if token != tokens[email] {
		data := atempt{
			Email: email,
			Error: "Invalid token",
		}

		err := m.render(w, data, "layout.html", "handle_code.html")
		if err != nil {
			m.httpError(w, err)

			return
		}

		return
	}

	delete(tokens, email)

	// Adding email to the context
	r = r.WithContext(context.WithValue(r.Context(), "email", email))
	m.afterLogin(w, r)
}
