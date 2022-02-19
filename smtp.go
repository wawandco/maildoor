package maildoor

import (
	"fmt"
	"net/smtp"
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
		msg := fmt.Sprintf("From: %v\n", opts.From)
		msg += fmt.Sprintf("To: %v\n", message.To)
		msg += fmt.Sprintf("Subject: %v\n", message.Subject)
		msg += "Content-Type: text/html\n\n"

		for _, b := range message.Bodies {
			if b.ContentType != "text/html" {
				continue
			}

			msg += string(message.Bodies[0].Content)
		}

		username := opts.From
		if opts.Username != "" {
			username = opts.Username
		}

		// Identity is empty usually.
		auth := smtp.PlainAuth("", username, opts.Password, opts.Host)
		err := smtp.SendMail(opts.Host+":"+opts.Port, auth, opts.From, []string{message.To}, []byte(msg))
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("error sending email: %w", err)
		}

		return nil
	}
}
