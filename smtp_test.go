package maildoor_test

import (
	"testing"

	smtpmock "github.com/mocktools/go-smtp-mock"
	"github.com/wawandco/maildoor"
)

func TestSMTPServer(t *testing.T) {
	server := smtpmock.New(smtpmock.ConfigurationAttr{
		LogToStdout:       true,
		LogServerActivity: true,
	})

	go func() {
		// To start server use Start() method
		if err := server.Start(); err != nil {
			t.Fatalf("error starting the mock smtp server: %s", err)
		}
	}()

	defer func() {
		err := server.Stop()
		if err != nil {
			panic(err)
		}
	}()

	// TODO: Start the server
	sender := maildoor.NewSMTPSender(maildoor.SMTPOptions{
		From:     "",
		Host:     "",
		Port:     "",
		Username: "",
		Password: "",
	})

	sender(&maildoor.Message{
		To:      "maildoor@wawand.co",
		Subject: "Login Message",
		Bodies: []*maildoor.MessageBody{
			{
				ContentType: "text/html",
				Content:     []byte("<h1>Hello World</h1>"),
			},
		},
	})

	// TODO: Check that the message was sent.

}
