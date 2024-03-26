package maildoor

import (
	"net/http"
)

// handleEmail endpoint validates the handleEmail and sends a token to the
// user by calling the handleEmail sender function.
func (m *maildoor) handleEmail(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	err := m.emailValidator(email)
	if err != nil {
		data := atempt{
			Error: err.Error(),
		}

		err := m.render(w, data, "layout.html", "handle_login.html")
		if err != nil {
			m.httpError(w, err)

			return
		}

		return
	}

	// token := newTokenFor(email, 6)
	data := atempt{}
	err = m.render(w, data, "layout.html", "handle_code.html")
	if err != nil {
		m.httpError(w, err)
		return
	}
}
