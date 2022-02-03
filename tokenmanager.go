package maildoor

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
