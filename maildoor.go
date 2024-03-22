package maildoor

import (
	"net/http"

	"github.com/wawandco/maildoor/internal"
)

// New maildoor handler
func New() http.Handler {
	s := http.NewServeMux()

	s.HandleFunc("GET /login/{$}", internal.Login)
	s.HandleFunc("POST /login/{$}", internal.Token)
	s.HandleFunc("POST /validate/{$}", internal.Validate)

	s.Handle("GET /", http.FileServer(http.FS(internal.Assets)))
	return s
}
