package maildoor

import (
	"context"
	"net/http"
)

// handleCode validates the input handleCode with the passed email.
func (m *maildoor) handleCode(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	code := r.FormValue("code")

	secret := secrets[email]

	if code == "" {
		renderError(m, w, email, "Invalid token")
		return
	}

	if !validateCode(code, secret) {
		renderError(m, w, email, "Invalid token")

		return
	}

	delete(secrets, email)

	// Adding email to the context
	r = r.WithContext(context.WithValue(r.Context(), "email", email))
	m.afterLogin(w, r)
}

func renderError(m *maildoor, w http.ResponseWriter, email, errorMessage string) {
	data := atempt{
		Email:       email,
		Error:       errorMessage,
		Logo:        m.logoURL,
		Icon:        m.iconURL,
		ProductName: m.productName,
	}

	err := m.render(w, data, "layout.html", "handle_code.html")
	if err != nil {
		m.httpError(w, err)

		return
	}
}
