package internal

import (
	"fmt"
	"html/template"
	"net/http"
	"sync"
)

var tux sync.Mutex
var tokens = map[string]string{}

func Token(prefix string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email := r.FormValue("email")

		tux.Lock()
		defer tux.Unlock()
		tokens[email] = "123456" // Token is generated here

		page := struct{
			Logo string
			Email string
			Error string
		}{
			Logo: "",
			Email: email,
		}
		tt := template.New("layout.html").Funcs(template.FuncMap{
			"prefixedPath": prefixedHelper(prefix),
		})

		// Generate a token and store it in the server
		// Save the user email in the session
		tt, err := tt.ParseFS(templates, "layout.html", "token.html")
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		err = tt.Execute(w, page)
		if err != nil {
			return
		}
	}
}
