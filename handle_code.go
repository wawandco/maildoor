package maildoor

import (
	"context"
	"net/http"
)

// handleCode validates the input handleCode with the passed email.
func (m *maildoor) handleCode(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	code := r.FormValue("code")

	// Find a combination of token and email in the server
	// call the afterlogin hook with the email
	// remove the token from the server
	if code != codes[email] {
		data := Attempt{
			Email:       email,
			Error:       "Invalid token",
			Logo:        m.logoURL,
			Icon:        m.iconURL,
			ProductName: m.productName,
		}

		html, err := m.codeRenderer(data)
		if err != nil {
			m.httpError(w, err)
			return
		}
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))

		return
	}

	delete(codes, email)

	// Adding email to the context
	r = r.WithContext(context.WithValue(r.Context(), "email", email))
	m.afterLogin(w, r)
}
