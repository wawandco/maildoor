package maildoor

import (
	"net/http"
)

// handleLogin enpoint renders the handleLogin page to enter the user
// identifier.
func (m *maildoor) handleLogout(w http.ResponseWriter, r *http.Request) {
	m.logout(w, r)
}
