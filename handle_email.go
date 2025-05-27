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
		
		html, renderErr := m.loginRenderer(data)
		if renderErr != nil {
			m.httpError(w, renderErr)
			return
		}
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
		return
	}

	token := newCodeFor(email)
	html, txt, err := m.mailBodies(token)
	if err != nil {
		m.httpError(w, err)
		return
	}

	err = m.emailSender(email, html, txt)
	if err != nil {
		data.Error = err.Error()
		w.WriteHeader(http.StatusInternalServerError)
		
		html, renderErr := m.loginRenderer(data)
		if renderErr != nil {
			m.httpError(w, renderErr)
			return
		}
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
		return
	}

	data.Email = email
	
	htmlContent, err := m.codeRenderer(data)
	if err != nil {
		m.httpError(w, err)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(htmlContent))
}
