package maildoor

import "net/http"

// option for the auth
type option func(*maildoor)

// Logo sets the logo url for the login page and the email
// that will be sent to the user.
func Logo(l string) option {
	return func(m *maildoor) {
		m.logoURL = l
	}
}

// ProductName allows to specify the product name used
// in emails and pages.
func ProductName(p string) option {
	return func(m *maildoor) {
		m.productName = p
	}
}

// Prefix sets the prefix for the routes. By default it is /auth/.
func Prefix(p string) option {
	return func(m *maildoor) {
		m.patternPrefix = p
	}
}

func Icon(i string) option {
	return func(m *maildoor) {
		m.iconURL = i
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

// Logout sets the function to be executed after logout
// this is useful to clear the session or cookie for the user
// and redirect to the login page. By default it redirects to
// the root of the app (/).
func Logout(fn func(http.ResponseWriter, *http.Request)) option {
	return func(m *maildoor) {
		m.logout = fn
	}
}

// EmailValidator sets the function to validate the email
// it can be replaced with a custom function.
func EmailValidator(fn func(email string) error) option {
	return func(m *maildoor) {
		m.emailValidator = fn
	}
}

// EmailSender is the function that will be called after the email
// has been determined to be valid. so the app can send the email to
// the user with the token. Txt and html are the email body in plain text and html format.
func EmailSender(fn func(to, html, txt string) error) option {
	return func(m *maildoor) {
		m.emailSender = fn
	}
}

// LoginRenderer sets a custom renderer function for the login page.
// The function receives the data needed to render the login page and should return
// the HTML string that will be sent to the user. If not set, the default template will be used.
func LoginRenderer(fn func(data Attempt) (string, error)) option {
	return func(m *maildoor) {
		m.loginRenderer = fn
	}
}

// CodeRenderer sets a custom renderer function for the code entry page.
// The function receives the data needed to render the code page and should return
// the HTML string that will be sent to the user. If not set, the default template will be used.
func CodeRenderer(fn func(data Attempt) (string, error)) option {
	return func(m *maildoor) {
		m.codeRenderer = fn
	}
}

// WithTokenStorage sets a custom token storage implementation.
// This allows you to use Redis, database, or any other storage backend
// instead of the default in-memory storage. The storage implementation
// must satisfy the ITokenStorage interface.
func WithTokenStorage(storage TokenStorage) option {
	return func(m *maildoor) {
		m.tokenStorage = storage
	}
}
