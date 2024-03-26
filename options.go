package maildoor

import "net/http"

// option for the auth
type option func(*maildoor)

// UsePrefix sets the prefix for the routes. By default it is /auth/.
func UsePrefix(p string) option {
	return func(m *maildoor) {
		m.patternPrefix = p
	}
}

// AfterLogin sets the function to be executed after login
// this is useful to set a cookie or a session for the user
// after the login is successful and redirect to secure area.
func AfterLogin(fn func(http.ResponseWriter, *http.Request)) option {
	return func(m *maildoor) {
		m.afterLogin = fn
	}
}

// EmailValidator sets the function to validate the email
// it can be replaced with a custom function.
func EmailValidator(fn func(email string) error) option {
	return func(m *maildoor) {
		m.emailValidator = fn
	}
}
