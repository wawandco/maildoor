package maildoor

import (
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// GenerateJWT token with the specified duration and secret.
func GenerateJWT(d time.Duration, secret []byte) (string, error) {
	expiration := time.Now().Add(d).Format(time.RFC3339)
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"ExpiresAt": expiration,
	})

	return t.SignedString(secret)
}

// ValidateJWT token with the specified secret.
func ValidateJWT(tt string, secret []byte) (bool, error) {
	tokenString := strings.TrimSpace(tt)
	t, err := jwt.ParseWithClaims(tokenString, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return secret, nil
	})

	if err != nil {
		return false, fmt.Errorf("error parsing error: %w", err)
	}

	cl := t.Claims.(*jwt.MapClaims)
	expires, err := time.Parse(time.RFC3339, (*cl)["ExpiresAt"].(string))

	if err != nil || expires.Before(time.Now()) {
		return false, nil
	}

	return err == nil, err
}
