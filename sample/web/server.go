package web

import (
	"embed"
	"fmt"
	"net/http"
	"os"

	"github.com/wawandco/maildoor"
)

var (
	//go:embed templates
	templates embed.FS

	secretKey = os.Getenv("SECRET_KEY")
)

// NewApp Builds the http Handler and in case of an error it returns it.
func NewApp() (http.Handler, error) {
	mux := http.NewServeMux()

	// Initialize the maildoor handler to take care of the web requests.
	auth, err := maildoor.NewWithOptions(
		secretKey,
		maildoor.UseTokenManager(maildoor.DefaultTokenManager(secretKey)),

		maildoor.UseFinder(finder),
		maildoor.UseSender(
			func(m *maildoor.Message) error {
				fmt.Println("Sending message: \n", string(m.Bodies[0].Content))

				return nil
			},
			// This could be a SMTP sender or other one.
			// maildoor.NewSMTPSender(maildoor.SMTPOptions{
			// 	From:     os.Getenv("SMTP_FROM_EMAIL"),
			// 	Host:     os.Getenv("SMTP_HOST"), // p.e. "smtp.gmail.com",
			// 	Port:     os.Getenv("SMTP_PORT"), //"587",
			// 	Password: os.Getenv("SMTP_PASSWORD"),
			// }),
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
