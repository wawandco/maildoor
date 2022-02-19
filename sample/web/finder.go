package web

import (
	"github.com/wawandco/maildoor"
	"github.com/wawandco/maildoor/sample"
)

func finder(token string) (maildoor.Emailable, error) {
	// maybe this needs to validate that the token is actually an email
	// for a simplistic example we just return a sample user.
	return sample.User(token), nil
}
