package sample

import (
	"bytes"
	_ "embed"
	"errors"
	"fmt"
	"net/http"
	"net/smtp"
	"os"
	"text/template"
	"time"

	"github.com/wawandco/maildoor"
)

// Auth handler with custom email validator
// and after login function
var Auth = maildoor.New(
	maildoor.UsePrefix("/auth/"),
	maildoor.WithLogo("https://raw.githubusercontent.com/wawandco/maildoor/5de0561/internal/sample/logo.png"),
	maildoor.ProductName("Basse"),
	maildoor.EmailValidator(validateEmail),
	maildoor.AfterLogin(afterLogin),
	maildoor.EmailSender(sendEmail),
	maildoor.Logout(logout),
)

// email struct to hold the email data to be used
// with the email template
type email struct {
	From    string
	To      string
	Subject string
	HTML    string
	Text    string
}

//go:embed email_template.txt
var mtmpl string

// emailTmpl is the template to be used to send the email
var emailTmpl = template.Must(template.New("email").Parse(mtmpl))

// sendEmail function to send the multipart email to the user
func sendEmail(to, html, txt string) error {
	from := os.Getenv("SMTP_FROM")
	password := os.Getenv("SMTP_PASS")
	user := os.Getenv("SMTP_USER")

	mb := bytes.NewBuffer([]byte{})
	err := emailTmpl.Execute(mb, email{
		HTML:    html,
		Text:    txt,
		From:    from,
		To:      to,
		Subject: "Your authentication code",
	})

	if err != nil {
		return fmt.Errorf("error executing email template: %w", err)
	}

	auth := smtp.PlainAuth("", user, password, "smtp.resend.com")
	err = smtp.SendMail("smtp.resend.com:587", auth, from, []string{to}, mb.Bytes())
	if err != nil {
		return fmt.Errorf("error sending smtp message: %w", err)
	}

	return nil
}

// afterLogin function to redirect the user to the private area
func afterLogin(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:    "sample",
		Value:   "sample",
		Expires: time.Now().Add(24 * time.Hour),
		Path:    "/",
	})

	http.Redirect(w, r, "/private", http.StatusFound)
}

// validateEmail function to validate the email address
// in this case we are only allowing a@pagano.id
func validateEmail(email string) error {
	if email == "a@pagano.id" {
		return nil
	}

	return errors.New("invalid email address")
}

// Logout clear the cookie and redirects to the root
func logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:    "sample",
		Value:   "",
		Expires: time.Now().Add(-1 * time.Hour),
		Path:    "/",
	})

	http.Redirect(w, r, "/", http.StatusFound)
}
