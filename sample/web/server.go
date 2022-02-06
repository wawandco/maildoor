package web

import (
	"embed"
	"fmt"
	"net/http"

	"github.com/wawandco/maildoor"
)

//go:embed templates
var templates embed.FS

func NewApp() (http.Handler, error) {
	mux := http.NewServeMux()

	// Initialize the maildoor handler to take care of the web requests.
	auth, err := maildoor.New(maildoor.Options{
		CSRFTokenSecret: "secret",
		SenderFn:        sender,

		FinderFn: finder,

		AfterLoginFn: afterLogin,
		LogoutFn:     logout,
	})

	if err != nil {
		return nil, fmt.Errorf("error initializing maildoor: %w", err)
	}

	mux.HandleFunc("/private", authenticated(private))
	mux.Handle("/auth/", auth)
	mux.HandleFunc("/", public)

	return mux, nil
}

func private(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	content, err := templates.ReadFile("templates/private.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.Write(content)
}

func public(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	content, err := templates.ReadFile("templates/public.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.Write(content)
}
