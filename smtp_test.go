package maildoor_test

import (
	"bufio"
	"bytes"
	"net"
	"net/textproto"
	"strings"
	"testing"

	"github.com/wawandco/maildoor"
	"github.com/wawandco/maildoor/testhelpers"
)

func runSMTPServer(t *testing.T) (func() string, string, error) {
	// prevent data race on bcmdbuf
	var done = make(chan struct{})
	var cmdbuf bytes.Buffer

	// The SMTP commants that the fake server returns
	serverData := []string{
		`220 hello world`,
		`502 EH?`,
		`250 mx.maildoor.server at your service`,
		`250 Sender ok`,
		`250 Receiver ok`,
		`354 Go ahead`,
		`250 Data ok`,
		`221 Goodbye`,
	}

	bcmdbuf := bufio.NewWriter(&cmdbuf)
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("Unable to create listener: %v", err)
	}

	go func() {
		defer close(done)

		conn, err := l.Accept()
		if err != nil {
			t.Errorf("Accept error: %v", err)
			return
		}
		defer conn.Close()

		tc := textproto.NewConn(conn)
		for i := 0; i < len(serverData) && serverData[i] != ""; i++ {
			tc.PrintfLine(serverData[i])
			for len(serverData[i]) >= 4 && serverData[i][3] == '-' {
				i++
				tc.PrintfLine(serverData[i])
			}

			if serverData[i] == "221 Goodbye" {
				return
			}

			read := false
			for !read || serverData[i] == "354 Go ahead" {
				msg, err := tc.ReadLine()
				bcmdbuf.Write([]byte(msg + "\r\n"))
				read = true
				if err != nil {
					t.Errorf("Read error: %v", err)
					return
				}
				if serverData[i] == "354 Go ahead" && msg == "." {
					break
				}
			}
		}
	}()

	closeFunc := func() string {
		<-done
		l.Close()

		bcmdbuf.Flush()
		return cmdbuf.String()
	}

	return closeFunc, l.Addr().String(), nil
}

func TestSMTPServer(t *testing.T) {
	closeFn, address, err := runSMTPServer(t)
	testhelpers.NoError(t, err)

	host := strings.Split(address, ":")[0]
	port := strings.Split(address, ":")[1]

	sender := maildoor.NewSMTPSender(maildoor.SMTPOptions{
		From:     "apagano@wawand.co",
		Host:     host,
		Port:     port,
		Username: "",
		Password: "",
	})

	err = sender(&maildoor.Message{
		To:      "maildoor@wawand.co",
		Subject: "Login Message",
		Bodies: []*maildoor.MessageBody{
			{
				ContentType: "text/html",
				Content:     []byte("<h1>Hello World</h1>"),
			},
			{
				ContentType: "text/plain",
				Content:     []byte("Hello World Plain"),
			},
		},
	})

	testhelpers.NoError(t, err)
	actualcmds := closeFn()

	testhelpers.Contains(t, actualcmds, "From: apagano@wawand.co")
	testhelpers.Contains(t, actualcmds, "To: maildoor@wawand.co")
	testhelpers.Contains(t, actualcmds, "Subject: Login Message")
	testhelpers.Contains(t, actualcmds, "Content-Type: text/html")
	testhelpers.Contains(t, actualcmds, "<h1>Hello World</h1>")

	// While we figure multipart email we'll only sent text/html
	testhelpers.NotContains(t, actualcmds, "Hello Plain")
	testhelpers.NotContains(t, actualcmds, "Content-Type: text/text")
}
