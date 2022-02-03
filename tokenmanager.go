package maildoor

import (
	"fmt"

	"github.com/dgrijalva/jwt-go"
)

// TokenManager interface allows to provide a custom token manager
// for maildoor to use. The default implementation of the tokenmanager
// is the JWT token manager which generates and validates tokens signed with
// JWT. This interface allows applications to replace JWT with other token
// validations mechanics such as Database or anything else.
type TokenManager interface {
	// Generate is expected to return a token string
	// that will be used as part of the email sent to then
	// be validated by the validate handler.
	Generate(user Emailable) (string, error)

	// Validate is the method in charge of validating a
	// received token, if the token is valid it should return true
	// if there is any error while validating (p.e. database connection)
	// it should return the error.
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
