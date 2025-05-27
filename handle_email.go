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
		
		// Use custom renderer if available
		if m.loginRenderer != nil {
			html, renderErr := m.loginRenderer(data)
			if renderErr != nil {
				m.httpError(w, renderErr)
				return
			}
			w.Header().Set("Content-Type", "text/html")
			w.Write([]byte(html))
			return
		}
		
		// Fall back to default template rendering
		err := m.render(w, data, "layout.html", "handle_login.html")
		if err != nil {
			m.httpError(w, err)
		}

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
		
		// Use custom renderer if available
		if m.loginRenderer != nil {
			html, renderErr := m.loginRenderer(data)
			if renderErr != nil {
				m.httpError(w, renderErr)
				return
			}
			w.Header().Set("Content-Type", "text/html")
			w.Write([]byte(html))
			return
		}
		
		// Fall back to default template rendering
		err := m.render(w, data, "layout.html", "handle_login.html")
		if err != nil {
			m.httpError(w, err)
		}

		return
	}

	data.Email = email
	
	// Use custom renderer if available
	if m.codeRenderer != nil {
		html, err := m.codeRenderer(data)
		if err != nil {
			m.httpError(w, err)
			return
		}
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
		return
	}
	
	// Fall back to default template rendering
	err = m.render(w, data, "layout.html", "handle_code.html")
	if err != nil {
		m.httpError(w, err)
		return
	}
}
