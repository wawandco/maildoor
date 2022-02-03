package maildoor

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// JWTTokenManager is the default token manager which is
// an alias for a byte slice that will implement the TokenManager
// interface by using JWT.
type JWTTokenManager []byte

// Generates a JWT token that lasts for 30 minutes. The duration of the token
// is specified within the ExpiresAt claim.
func (dm JWTTokenManager) Generate(user Emailable) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Minute * 30).Unix(),
	})

	return token.SignedString([]byte(dm))
}

// Validate the passed token and returns true if the token is valid. This implementation
// checks that the ExpiresAt claim to check that the token has not expired.
func (dm JWTTokenManager) Validate(tt string) (bool, error) {
	t, err := jwt.Parse(tt, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(dm), nil
	})

	cl := t.Claims.(jwt.StandardClaims)
	if cl.ExpiresAt < time.Now().Unix() {
		return false, fmt.Errorf("token expired")
	}

	return err == nil, err
}
