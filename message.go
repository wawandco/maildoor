package maildoor

type Message struct {
	To      string
	Subject string
	Bodies  []*Body
}

func (m *Message) AddBody(contentType string, content []byte) {
	m.Bodies = append(m.Bodies, &Body{
		ContentType: contentType,
		Content:     content,
	})
}

type Body struct {
	ContentType string
	Content     []byte
}
