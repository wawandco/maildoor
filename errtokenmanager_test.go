package maildoor_test

import (
	"fmt"

	"github.com/wawandco/maildoor"
)

type errTokenManager string

func (et errTokenManager) Generate(maildoor.Emailable) (string, error) {
	return "", fmt.Errorf("%s", string(et))
}

func (et errTokenManager) Validate(tt string) (bool, error) {
	return true, fmt.Errorf("%s", string(et))
}
