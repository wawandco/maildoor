package maildoor

import (
	"net/http"
)

// handleLogin enpoint renders the handleLogin page to enter the user
// identifier.
func (m *maildoor) handleLogin(w http.ResponseWriter, r *http.Request) {
	err := m.render(w, atempt{}, "layout.html", "handle_login.html")
	if err != nil {
		m.httpError(w, err)

		return
	}
}
