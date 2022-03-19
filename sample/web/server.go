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

// NewApp Builds the http Handler and in case of an error it returns it.
func NewApp() (http.Handler, error) {
	mux := http.NewServeMux()

	// Initialize the maildoor handler to take care of the web requests.
	auth, err := maildoor.NewWithOptions(
		os.Getenv("SECRET_KEY"),

		maildoor.UseFinder(finder),
		maildoor.UseAfterLogin(afterLogin),
		maildoor.UseLogout(logout),
		maildoor.UseTokenManager(maildoor.DefaultTokenManager(os.Getenv("SECRET_KEY"))),
		maildoor.UseSender(
			maildoor.NewSMTPSender(maildoor.SMTPOptions{
				From:     os.Getenv("SMTP_FROM_EMAIL"),
				Host:     os.Getenv("SMTP_HOST"), // p.e. "smtp.gmail.com",
				Port:     os.Getenv("SMTP_PORT"), //"587",
				Password: os.Getenv("SMTP_PASSWORD"),
			}),
		),
	)

	if err != nil {
		return nil, fmt.Errorf("error initializing maildoor: %w", err)
	}

	mux.HandleFunc("/private", authenticated(private))
	mux.Handle("/auth/", auth)
	mux.HandleFunc("/", public)

	return mux, nil
}
