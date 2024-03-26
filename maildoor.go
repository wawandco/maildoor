package maildoor

import (
	"embed"
	"html/template"
	"log/slog"
	"net/http"
	"path"
	"strings"
	"sync"
	"time"
)

var tux sync.Mutex
var tokens = map[string]string{}

var (
	//go:embed *.html *.txt
	templates embed.FS

	//go:embed *.png
	assets embed.FS
)

// New maildoor handler with the passed options.
func New(options ...option) http.Handler {
	s := &maildoor{
		mux: http.NewServeMux(),

		afterLogin: func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Logged in!"))
		},

		emailValidator: func(email string) error {
			// All emails are valid by default
			return nil
		},
	}

	for _, opt := range options {
		opt(s)
	}

	s.HandleFunc("GET /login", s.handleLogin)
	s.HandleFunc("POST /email", s.handleEmail)
	s.HandleFunc("POST /code", s.handleCode)

	// Adding the static assets handler
	ah := http.StripPrefix(s.patternPrefix, http.FileServer(http.FS(assets)))
	s.Handle("GET /*", ah)

	return s
}

type maildoor struct {
	mux *http.ServeMux

	patternPrefix  string
	afterLogin     http.HandlerFunc
	emailValidator func(email string) error
}

func (m *maildoor) HandleFunc(pattern string, handler http.HandlerFunc) {
	// prefix the pattens with the routesPrefix
	parts := strings.Split(pattern, " ")
	pattern = path.Join(m.patternPrefix, parts[0])
	if len(parts) == 2 {
		pattern = path.Join(m.patternPrefix, parts[1])
		pattern = parts[0] + " " + pattern
	}

	// Adding the handler to the ServeMux
	m.mux.HandleFunc(pattern, handler)
}

func (m *maildoor) Handle(pattern string, handler http.Handler) {
	// prefix the pattens with the routesPrefix
	parts := strings.Split(pattern, " ")
	pattern = path.Join(m.patternPrefix, parts[0])
	if len(parts) == 2 {
		pattern = path.Join(m.patternPrefix, parts[1])
		pattern = parts[0] + " " + pattern
	}

	// Adding the handler to the ServeMux
	m.mux.Handle(pattern, handler)
}

func (m *maildoor) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Adding common things here, loggers and other things.
	t := time.Now()

	// Parsing form
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	m.mux.ServeHTTP(w, r)
	slog.Info(">", "method", r.Method, "path", r.URL.Path, "duration", time.Since(t))
}

func (m *maildoor) helpers() template.FuncMap {
	return template.FuncMap{
		"prefixedPath": func(p string) string {
			return path.Join(m.patternPrefix, p)
		},
	}
}
