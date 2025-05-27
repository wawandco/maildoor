package maildoor

import (
	"bytes"
	"embed"
	"html/template"
	"io"
	"log/slog"
	"net/http"
	"path"
	"strings"
	"time"
)

var (
	//go:embed *.html *.txt
	templates embed.FS

	//go:embed *.png
	assets embed.FS
)

// Attempt is a struct to hold the email and error message.
// used across different views.
type Attempt struct {
	Logo        string
	Icon        string
	ProductName string
	Email       string
	Error       string
	Code        string
}

// New maildoor handler with the passed options.
func New(options ...option) http.Handler {
	s := &maildoor{
		mux:         http.NewServeMux(),
		productName: "Maildoor",
		logoURL:     "https://raw.githubusercontent.com/wawandco/maildoor/508ff43/assets/images/maildoor_logo.png",
		iconURL:     "https://raw.githubusercontent.com/wawandco/maildoor/508ff43/assets/images/maildoor_icon.png",

		afterLogin: func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Logged in!"))
		},

		emailValidator: func(email string) error {
			// All emails are valid by default
			return nil
		},

		logout: func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "/", http.StatusFound)
		},

	}

	// Set default login renderer
	s.loginRenderer = s.defaultLoginRenderer
	
	// Set default code renderer
	s.codeRenderer = s.defaultCodeRenderer

	for _, opt := range options {
		opt(s)
	}

	s.HandleFunc("GET /login", s.handleLogin)
	s.HandleFunc("POST /email", s.handleEmail)
	s.HandleFunc("POST /code", s.handleCode)
	s.HandleFunc("DELETE /logout", s.handleLogout)

	// Adding the static assets handler
	ah := http.StripPrefix(s.patternPrefix, http.FileServer(http.FS(assets)))
	s.Handle("GET /*", ah)

	return s
}

type maildoor struct {
	mux *http.ServeMux

	productName string
	logoURL     string
	iconURL     string

	patternPrefix string
	afterLogin    http.HandlerFunc
	logout        http.HandlerFunc

	emailValidator func(email string) error
	emailSender    func(email, html, txt string) error

	loginRenderer func(data Attempt) (string, error)
	codeRenderer  func(data Attempt) (string, error)
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
		m.httpError(w, err)
		return
	}

	// Correcting method based on _method field
	if r.Method == "POST" && r.FormValue("_method") != "" {
		r.Method = r.FormValue("_method")
	}

	m.mux.ServeHTTP(w, r)
	slog.Info(">", "method", r.Method, "path", r.URL.Path, "duration", time.Since(t))
}

// render a template with the passed data and partials using
// the templates FS. if using layout it should go first.
func (m *maildoor) render(w io.Writer, data any, partials ...string) error {
	if len(partials) == 0 {
		return nil
	}

	tt := template.New(partials[0]).Funcs(template.FuncMap{
		"prefixedPath": func(p string) string {
			return path.Join(m.patternPrefix, p)
		},
	})

	tt, err := tt.ParseFS(templates, partials...)
	if err != nil {
		return err
	}

	return tt.Execute(w, data)
}

func (m *maildoor) httpError(w http.ResponseWriter, err error) {
	slog.Error("*", "error", err.Error())
	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
}

// defaultLoginRenderer is the default renderer for login pages
func (m *maildoor) defaultLoginRenderer(data Attempt) (string, error) {
	var buf bytes.Buffer
	err := m.render(&buf, data, "layout.html", "handle_login.html")
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

// defaultCodeRenderer is the default renderer for code pages
func (m *maildoor) defaultCodeRenderer(data Attempt) (string, error) {
	var buf bytes.Buffer
	err := m.render(&buf, data, "layout.html", "handle_code.html")
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (m *maildoor) mailBodies(code string) (string, string, error) {
	data := struct {
		Code    string
		Logo    string
		Product string
		Year    string
	}{
		Code:    code,
		Logo:    m.logoURL,
		Product: m.productName,
		Year:    time.Now().Format("2006"),
	}

	sw := bytes.NewBuffer([]byte{})
	err := m.render(sw, data, "message.html")
	if err != nil {
		return "", "", err
	}

	html := sw.String()

	sw = bytes.NewBuffer([]byte{})
	err = m.render(sw, data, "message.txt")
	if err != nil {
		return "", "", err
	}

	txt := sw.String()

	return html, txt, nil
}
