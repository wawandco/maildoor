package maildoor

import (
	"html/template"
	"math/rand"
	"net/http"
)

// handleEmail endpoint validates the handleEmail and sends a token to the
// user by calling the handleEmail sender function.
func (m *maildoor) handleEmail(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	err := m.emailValidator(email)
	if err != nil {
		msg := err.Error()

		// Render login again with error.
		tt := template.New("layout.html").Funcs(m.helpers())
		tt, err = tt.ParseFS(templates, "layout.html", "handle_login.html")
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		err = tt.Execute(w, struct {
			Error string
		}{
			Error: msg,
		})

		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		return
	}

	// Generating a new token
	var letters = []rune("1234567890")
	b := make([]rune, 6)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	token := string(b)

	tux.Lock()
	defer tux.Unlock()
	tokens[email] = token

	// Generate a token and store it in the server
	// Save the user email in the session
	tt := template.New("layout.html").Funcs(m.helpers())
	tt, err = tt.ParseFS(templates, "layout.html", "handle_code.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = tt.Execute(w, struct {
		Email string
		Error string
	}{
		Email: email,
	})

	if err != nil {
		return
	}
}
