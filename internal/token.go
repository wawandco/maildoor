package internal

import (
	"html/template"
	"math/rand"
	"net/http"
	"sync"
)

var tux sync.Mutex
var tokens = map[string]string{}

// NewToken generates a numeric token of 6 digits.
// to be used as a temporary password.
func newToken() string {
	var letters = []rune("1234567890")

    b := make([]rune, 6)
    for i := range b {
        b[i] = letters[rand.Intn(len(letters))]
    }

    return string(b)
}

func Token(prefix string, emailValidator func(string)(error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tt := template.New("layout.html").Funcs(template.FuncMap{
			"prefixedPath": prefixedHelper(prefix),
		})

		email := r.FormValue("email")
		err := emailValidator(email)
		if err != nil {
			msg := err.Error()

			// Render login again with error.
			tt, err = tt.ParseFS(templates, "layout.html", "login.html")
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			err = tt.Execute(w, struct{
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

		tux.Lock()
		defer tux.Unlock()
		token := newToken()
		tokens[email] = token

		// Generate a token and store it in the server
		// Save the user email in the session
		tt, err = tt.ParseFS(templates, "layout.html", "token.html")
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		err = tt.Execute(w, struct{
			Email string
			Error string
		}{
			Email: email,
		})

		if err != nil {
			return
		}
	}
}
