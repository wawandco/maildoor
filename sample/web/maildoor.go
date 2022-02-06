package web

import (
	"fmt"
	"net/http"
	"net/smtp"
	"os"
	"time"

	"github.com/wawandco/maildoor"
	"github.com/wawandco/maildoor/sample"
)

func afterLogin(w http.ResponseWriter, r *http.Request, user maildoor.Emailable) error {
	// Sets the sample cookie so the user can pass the
	// authenticated middleware.
	cookie := &http.Cookie{
		Name:     "sample",
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
		Secure:   true,
		Value:    user.EmailAddress(),
	}

	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/private", http.StatusSeeOther)

	return nil
}

// This is a sample Sendgrid SMTP sender, DO NOT USE THIS in production.
// It reads from the SENDGRID_API_KEY and SENDGRID_FROM_EMAIL variables
// to send emails. These must be set for the sender to work correctly.
//
// Caveats:
// - It only uses the first body in the message (text/html)
// - Subject is hardcoded.
// - It does not handle empty environment variables.
func sender(message *maildoor.Message) error {
	from := os.Getenv("SENDGRID_FROM_EMAIL")

	msg := fmt.Sprintf("From: %v\n", from)
	msg += fmt.Sprintf("To: %v\n", message.To)
	msg += fmt.Sprintf("Subject: %v\n", "Your login link to sample app")
	msg += "Content-Type: text/html\n\n"
	msg += string(message.Bodies[0].Content)

	auth := smtp.PlainAuth("", "apikey", os.Getenv("SENDGRID_API_KEY"), "smtp.sendgrid.net")
	err := smtp.SendMail("smtp.sendgrid.net:587", auth, from, []string{message.To}, []byte(msg))
	if err != nil {
		return fmt.Errorf("error sending email: %w", err)
	}

	return nil
}

func finder(token string) (maildoor.Emailable, error) {
	// maybe this needs to validate that the token is actually an email
	// for a simplistic example we just return a sample user.
	return sample.User(token), nil
}

// logout function clears the sample cookie.
func logout(w http.ResponseWriter, r *http.Request) error {
	cookie := &http.Cookie{
		Name:     "sample",
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Now().Add(-1),
		Secure:   true,

		// This will expire the cookie.
		MaxAge: -1,
	}

	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/", http.StatusSeeOther)

	return nil
}
