package web

import (
	"embed"
	"fmt"
	"net/http"
	"os"

	"github.com/wawandco/maildoor"
)

//go:embed templates
var templates embed.FS

func NewApp() (http.Handler, error) {
	mux := http.NewServeMux()

	// Initialize the maildoor handler to take care of the web requests.
	auth, err := maildoor.New(maildoor.Options{
		CSRFTokenSecret: "secret",
		FinderFn:        finder,

		SenderFn: maildoor.NewSMTPSender(maildoor.SMTPOptions{
			From:     os.Getenv("SMTP_FROM_EMAIL"),
			Host:     os.Getenv("SMTP_HOST"), //"smtp.sendgrid.net",
			Port:     os.Getenv("SMTP_PORT"), //"587",
			Username: os.Getenv("SMTP_USERNAME"),
			Password: os.Getenv("STP_PASSWORD"),
		}),

		AfterLoginFn: afterLogin,
		LogoutFn:     logout,

		// TokenManager using the secret key
		TokenManager: maildoor.DefaultTokenManager(os.Getenv("TOKEN_MANAGER_SECRET")),
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
