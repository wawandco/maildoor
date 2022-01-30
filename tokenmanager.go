package maildoor

import (
	"fmt"

	"github.com/dgrijalva/jwt-go"
)

// TokenManager interface allows to provide a custom token manager
// for maildoor to use.
type TokenManager interface {
	Generate(user Emailable) (string, error)
	Validate(token string) (bool, error)
}

// JWTTokenManager is the default token manager which is
// an alias for a byte slice that will implement the TokenManager
// interface by using JWT.
type JWTTokenManager []byte

func (dm JWTTokenManager) Generate(user Emailable) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{})
	return token.SignedString([]byte(dm))
}

func (dm JWTTokenManager) Validate(tt string) (bool, error) {
	_, err := jwt.Parse(tt, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(dm), nil
	})

	return err == nil, err
}
