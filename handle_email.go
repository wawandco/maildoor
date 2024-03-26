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

	token := newTokenFor(email)
	html, txt, err := m.mailBodies(token)
	if err != nil {
		m.httpError(w, err)
		return
	}

	err = m.emailSender(email, html, txt)
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

	data := atempt{
		Email: email,
	}

	err = m.render(w, data, "layout.html", "handle_code.html")
	if err != nil {
		m.httpError(w, err)
		return
	}
}
