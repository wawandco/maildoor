package maildoor

import "time"

// DefaultTokenManager is the default token manager which is
// an alias for a byte slice that will implement the TokenManager
// interface by using JWT.
type DefaultTokenManager []byte

// Generates a JWT token that lasts for 30 minutes. The duration of the token
// is specified within the ExpiresAt claim.
func (dm DefaultTokenManager) Generate(user Emailable) (string, error) {
	return GenerateJWT(30*time.Minute, []byte(dm))
}

// Validate the passed token and returns true if the token is valid. This implementation
// checks that the ExpiresAt claim to check that the token has not expired.
func (dm DefaultTokenManager) Validate(token string) (bool, error) {
	return ValidateJWT(token, []byte(dm))
}
