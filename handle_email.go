package maildoor

import (
	"net/http"
)

// handleEmail endpoint validates the handleEmail and sends a token to the
// user by calling the handleEmail sender function.
func (m *maildoor) handleEmail(w http.ResponseWriter, r *http.Request) {
	data := Attempt{
		Logo:        m.logoURL,
		ProductName: m.productName,
		Icon:        m.iconURL,
	}

	email := r.FormValue("email")
	if err := m.emailValidator(email); err != nil {
		data.Error = err.Error()
		w.WriteHeader(http.StatusUnprocessableEntity)

		html, err := m.loginRenderer(data)
		if err != nil {
			m.httpError(w, err)
			return
		}

		w.Header().Set("Content-Type", "text/html")
		_, err = w.Write([]byte(html))
		if err != nil {
			m.httpError(w, err)
			return
		}

		return
	}

	token := m.newCodeFor(email)
	html, txt, err := m.mailBodies(token)
	if err != nil {
		m.httpError(w, err)
		return
	}

	err = m.emailSender(email, html, txt)
	if err != nil {
		data.Error = err.Error()
		w.WriteHeader(http.StatusInternalServerError)

		html, err := m.loginRenderer(data)
		if err != nil {
			m.httpError(w, err)
			return
		}

		w.Header().Set("Content-Type", "text/html")
		_, err = w.Write([]byte(html))
		if err != nil {
			m.httpError(w, err)
			return
		}

		return
	}

	data.Email = email

	htmlContent, err := m.codeRenderer(data)
	if err != nil {
		m.httpError(w, err)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	_, err = w.Write([]byte(htmlContent))
	if err != nil {
		m.httpError(w, err)
		return
	}
}
