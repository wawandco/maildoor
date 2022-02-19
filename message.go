package maildoor

// Message allows to pass the message content to the sender,
// that way the sender function can take care of sending the
// message by using the most appropriate email sending provider
// for the app.
type Message struct {
	From    string
	To      string
	Subject string
	Bodies  []*MessageBody
}

func (m *Message) addBody(contentType string, content []byte) {
	m.Bodies = append(m.Bodies, &MessageBody{
		ContentType: contentType,
		Content:     content,
	})
}

// MessageBody is a part of the email message, maildoor
// adds text/plain and text/html to the message, for each of these
// bodies it adds respective content.
type MessageBody struct {
	ContentType string
	Content     []byte
}
