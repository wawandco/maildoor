package maildoor

// Emailable is the type that will be returned from the
// finder, for maildoor to work it needs to be able to
// use a type that can provide an email address.
type Emailable interface {
	// EmailAddress returns the email address of the user
	// that will be authenticated. This address is used to
	// send the authentication email to the user.
	EmailAddress() string
}
