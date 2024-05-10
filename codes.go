package maildoor

import (
	"math/rand"
	"sync"
)

var (
	tux     sync.Mutex
	codes   = map[string]string{}
	letters = []rune("1234567890")
)

// newCodeFor generates a new code for the email and stores it in the codes map.
// tokens are always 6 characters long.
func newCodeFor(email string) string {
	// Generating a new token
	b := make([]rune, 6)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	tux.Lock()
	defer tux.Unlock()
	codes[email] = string(b)

	return string(b)
}
