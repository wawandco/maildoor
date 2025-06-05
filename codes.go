package maildoor

import (
	"math/rand"
)

var (
	letters = []rune("1234567890")
)

// newCodeFor generates a new code for the email and stores it using the TokenStorage.
// tokens are always 6 characters long.
func (m *maildoor) newCodeFor(email string) string {
	// Generating a new token
	b := make([]rune, 6)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	token := string(b)
	m.tokenStorage.Store(email, token)

	return token
}
