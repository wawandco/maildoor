package maildoor

import (
	"net/http"
)

// handleLogin enpoint renders the handleLogin page to enter the user
// identifier.
func (m *maildoor) handleLogin(w http.ResponseWriter, r *http.Request) {
	data := Attempt{
		Logo:        m.logoURL,
		Icon:        m.iconURL,
		ProductName: m.productName,
	}

	html, err := m.loginRenderer(data)
	if err != nil {
		m.httpError(w, err)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}
