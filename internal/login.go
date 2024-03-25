package internal

import (
	"fmt"
	"html/template"
	"net/http"
)

// Login enpoint renders the login page to enter the user
// identifier.
func Login(prefix string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tt := template.New("layout.html").Funcs(template.FuncMap{
			"prefixedPath": prefixedHelper(prefix),
		})

		tt, err := tt.ParseFS(templates, "layout.html", "login.html")
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		err = tt.Execute(w, struct{
			Error string
		}{})

		if err != nil {
			return
		}
	}
}
