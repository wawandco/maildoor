package maildoor

import "net/http"

// option for the auth
type option func()

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
