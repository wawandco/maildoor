package maildoor

import (
	"fmt"
	"net/smtp"
	"time"
)

// SMTPOptions are the options for the SMTP sender.
type SMTPOptions struct {
	// The email address of the sender, this is the one used in
	// the message.
	From string

	// Host to authenticate/send the email.
	Host string

	// Port to authenticate/send the email.
	Port string

	// Username to authenticate/send the email, if not specified, the From field
	// is used.
	Username string

	// Password to authenticate/send the email. its recommended to use an environment
	// variable to pull this.
	Password string
}

//NewSMTPSender creates a new sender that will use the options
//to send the message HTML through SMTP.
func NewSMTPSender(opts SMTPOptions) func(*Message) error {
	return func(message *Message) error {
		message.From = opts.From

		username := opts.From
		if opts.Username != "" {
			username = opts.Username
		}

		// Identity is empty usually.
		auth := smtp.PlainAuth("", username, opts.Password, opts.Host)
		err := smtp.SendMail(opts.Host+":"+opts.Port, auth, opts.From, []string{message.To}, smtpContent(message))
		if err != nil {
			return fmt.Errorf("error sending email: %w", err)
		}

		return nil
	}
}

// Builds the content for the multipart message.
func smtpContent(m *Message) []byte {
	boundary := time.Now().UnixNano()

	msg := fmt.Sprintf("From: %v\n", m.From)
	msg += fmt.Sprintf("To: %v\n", m.To)
	msg += fmt.Sprintf("Subject: %v\n", m.Subject)
	msg += fmt.Sprintf("Content-Type: multipart/alternative; boundary=\"%x\"\n\n", boundary)

	for i, b := range m.Bodies {
		msg += fmt.Sprintf("--%x\n", boundary)
		msg += fmt.Sprintf("Content-Type: %s; charset=\"UTF-8\"\n\n", b.ContentType)

		msg += string(b.Content)
		msg += "\n"
		if i != len(m.Bodies)-1 {
			continue
		}

		msg += fmt.Sprintf("--%x--", boundary)
	}

	return []byte(msg)
}
