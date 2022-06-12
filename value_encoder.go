package maildoor

import (
	"fmt"
	"strings"

	"github.com/golang-jwt/jwt/v4"
)

type valueEncoder string

func (v valueEncoder) Encode(value string) (string, error) {
	t := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"Value": value,
		},
	)

	vx, err := t.SignedString([]byte(v))
	return vx, err
}

func (v valueEncoder) Decode(tt string) (string, error) {
	tokenString := strings.TrimSpace(tt)
	t, err := jwt.ParseWithClaims(
		tokenString,
		&jwt.MapClaims{},
		jwtKeyFunc([]byte(v)),
	)

	if err != nil {
		return "", fmt.Errorf("error parsing: %w", err)
	}

	cl, ok := t.Claims.(*jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("error matching claims: not a map claims instance")
	}

	return (*cl)["Value"].(string), nil
}
