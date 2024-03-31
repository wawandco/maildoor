![tests workflow](https://github.com/wawandco/maildoor/actions/workflows/test.yml/badge.svg)
![report card](https://goreportcard.com/badge/github.com/wawandco/maildoor)

# Maildoor

Maildoor is an email based authentication library that allows users to sign up and sign in to your application using their email address. It is a pluggable library that can be used with any go http server.

### Usage

Using maildoor is as simple as creating a new instance of the maildoor.Handler and passing it to your http server.

```go
// Initialize the maildoor handler
auth := maildoor.New(
	maildoor.WithLogo("https://example.com/logo.png"),
	maildoor.ProductName("My App"))
	maildoor.UsePrefix("/auth/"), // Prefix for the routes

	maildoor.EmailSender(func(email string, token string) error {
		// send email
		return nil
	}),
	maildoor.EmailValidator(func(email string) bool {
		// validate email
		return true
	}),
})

mux := http.NewServeMux()
mux.Handle("/auth", auth)
mux.Handle("/private", secure(privateHandler))
http.ListenAndServe(":8080", mux)
```

## Features

- Pluggable http.Handler that can be used with any go http server
- Customizable email sending mechanism
- Customizable email validation mechanism
- Customizable logo
- Customizable product name

### Roadmap

- Customizable templates (Bring your own).
- Custom token storage mechanism
- Time based token expiration out the box
- Prevend CSRF attacks with token
