package maildoor

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func GenerateJWT(d time.Duration, secret []byte) (string, error) {
	expiration := time.Now().Add(d).Format(time.RFC3339)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"ExpiresAt": expiration,
	})

	return token.SignedString(secret)
}

func ValidateJWT(tt string, secret []byte) (bool, error) {
	t, err := jwt.Parse(tt, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return secret, nil
	})

	if err != nil {
		return false, fmt.Errorf("error parsing error: %w", err)
	}

	cl := t.Claims.(jwt.MapClaims)
	expires, err := time.Parse(time.RFC3339, cl["ExpiresAt"].(string))

	if err != nil || expires.Before(time.Now()) {
		return false, fmt.Errorf("token expired")
	}

	return err == nil, err
}
