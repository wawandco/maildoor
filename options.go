package maildoor

import "net/http"

var (
	// routesPrefix for the maildoor routes
	routesPrefix = "/"

	// afterLogin is the function to be executed after login
	// this is useful to set a cookie or a session for the user
	afterLogin = func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Logged in!"))
	}

	// emailValidator is the function to validate the email
	// it can be replaced with a custom function.
	emailValidator = func(email string) (error) {
		return nil
	}
)

// option for the auth
type option func()

// UsePrefix sets the prefix for the routes. By default it is /auth/.
func UsePrefix(p string) option {
	return func() {
		routesPrefix = p
	}
}

// AfterLogin sets the function to be executed after login
// this is useful to set a cookie or a session for the user
// after the login is successful and redirect to secure area.
func AfterLogin(fn func(http.ResponseWriter, *http.Request)) option {
	return func() {
		afterLogin = fn
	}
}

// EmailValidator sets the function to validate the email
// it can be replaced with a custom function.
func EmailValidator(fn func(email string) (error)) option {
	return func() {
		emailValidator = fn
	}
}
