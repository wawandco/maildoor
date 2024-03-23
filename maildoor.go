package maildoor

import (
	"net/http"
	"path"

	"github.com/wawandco/maildoor/internal"
)

// routesPrefix for the maildoor routes
var routesPrefix = "/"
var afterLogin = func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Logged in!"))
}

// New maildoor handler
func New(options ...option) http.Handler {
	// Applying options before creating the handler
	for _, opt := range options {
		opt()
	}

	// withPrefix is a helper function to add the prefix to the routes
	withPrefix := func(p string) string {
		return path.Join(routesPrefix,p)
	}

	s := http.NewServeMux()
	s.HandleFunc("GET "+withPrefix("/login"), internal.Login(routesPrefix))
	s.HandleFunc("POST "+withPrefix("/login"), internal.Token(routesPrefix))
	s.HandleFunc("POST "+withPrefix("/validate"), internal.Validate(routesPrefix, afterLogin))

	// Adding the static assets handler
	s.Handle(
		"GET "+withPrefix("/*"),
		http.StripPrefix(routesPrefix, http.FileServer(http.FS(internal.Assets))),
	)

	return s
}
