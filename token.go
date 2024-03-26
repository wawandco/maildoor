package maildoor

import (
	"math/rand"
	"sync"
)

var (
	tux     sync.Mutex
	tokens  = map[string]string{}
	letters = []rune("1234567890")
)

// newTokenFor generates a new token for the email and stores it in the tokens map.
// tokens are always 6 characters long.
func newTokenFor(email string) string {
	// Generating a new token
	b := make([]rune, 6)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	tux.Lock()
	defer tux.Unlock()
	tokens[email] = string(b)

	return string(b)
}
