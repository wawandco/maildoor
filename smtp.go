package maildoor

import (
	"fmt"
	"net/smtp"
)

// SMTPOptions are the options for the SMTP sender.
type SMTPOptions struct {
	From string
	Host string
	Port string

	Username string
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

		auth := smtp.PlainAuth("", opts.Username, opts.Password, opts.Host)
		err := smtp.SendMail(opts.Host+":"+opts.Port, auth, opts.From, []string{message.To}, []byte(msg))
		if err != nil {
			return fmt.Errorf("error sending email: %w", err)
		}

		return nil
	}
}
