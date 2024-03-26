package maildoor

import (
	"html/template"
	"net/http"
)

// handleLogin enpoint renders the handleLogin page to enter the user
// identifier.
func (m *maildoor) handleLogin(w http.ResponseWriter, r *http.Request) {
	tt := template.New("layout.html").Funcs(m.helpers())
	tt, err := tt.ParseFS(templates, "layout.html", "handle_login.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = tt.Execute(w, struct {
		Error string
	}{})

	if err != nil {
		return
	}
}
