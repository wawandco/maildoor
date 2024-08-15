package maildoor

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"fmt"
	"sync"
	"time"
)

var (
	tux     sync.Mutex
	secrets = map[string]string{}
)

const (
	// expiresTime is the time in seconds that the code expires
	// 2 minutes
	expiresTime = 120
)

// generateSecret generates a 16 byte base32 encoded secret
// The secret is then returned
func generateSecret() (string, error) {
	secret := make([]byte, 10)
	_, err := rand.Read(secret)
	if err != nil {
		return "", err
	}

	return base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(secret), nil
}

// generateCode calls the gen function to generate a 6-digit code
// using the secret and the time step
func generateCode(secret string) string {
	timeStep := time.Now().Unix() / expiresTime
	return gen(secret, timeStep)
}

// validateCode validates the codeToValid with the secret
// by generating the code for the current time step and the previous and next time steps
// If the code matches any of the generated codes, it returns true otherwise false
func validateCode(codeToValid, secret string) bool {
	timeStep := time.Now().Unix() / expiresTime
	// Check the current time step, the previous and the next time steps
	// to allow for some time drift
	for i := -1; i <= 1; i++ {
		code := gen(secret, timeStep+int64(i))
		if code == codeToValid {
			return true
		}
	}

	return false
}

// gen generates a 6-digit code using the secret and the time step
// The code is then returned
func gen(secret string, timeStep int64) string {
	timeBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(timeBytes, uint64(timeStep))

	key := []byte(secret)
	hmacSha1 := hmac.New(sha1.New, key)
	hmacSha1.Write(timeBytes)
	hash := hmacSha1.Sum(nil)

	offset := hash[len(hash)-1] & 0x0f
	codeBytes := binary.BigEndian.Uint32(hash[offset : offset+4])

	code := codeBytes % 1_000_000
	codeStr := fmt.Sprintf("%06d", code)

	return codeStr
}

func saveSecret(email, secret string) {
	tux.Lock()
	defer tux.Unlock()
	secrets[email] = secret
}
