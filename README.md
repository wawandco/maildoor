![report card](https://goreportcard.com/badge/github.com/wawandco/maildoor)

# Maildoor

Maildoor is an email based authentication library that allows users to sign up and sign in to your application using their email address. It is a pluggable library that can be used with any go http server.

### Usage

Using maildoor is as simple as creating a new instance of the maildoor.Handler and passing it to your http server.

```go
// Initialize the maildoor handler
auth := maildoor.New(
	maildoor.Logo("https://example.com/logo.png"),
	maildoor.ProductName("My App"))
	maildoor.Prefix("/auth/"), // Prefix for the routes

	// Defines the email sending mechanism which is up to the
	// host application to implement.
	maildoor.EmailSender(func(to, html, txt string) error{
		// Send email to the user that's loggin in'
		return smtp.Send(to, html, txt)
	}),

	// Defines the email validation mechanism
	maildoor.EmailValidator(func(email string) bool {
		// Validate email with the users package
		return users.UserExists(email)
	}),

	// Defines what to do after the user has successfuly logged in
	// This is where you would set the user session or redirect to a private page
	maildoor.AfterLogin(func w http.ResponseWriter, r http.Request) {
		// Redirect to the private page
		http.Redirect(w, r, "/private", http.StatusFound)
	}),

	// Defines what to do after the user has successfuly loged out
	// This is where you would clear the user session or redirect to a login page
	maildoor.Logout(func(w http.ResponseWriter, r *http.Request){
		http.Redirect(w, r, "/auth/login", http.StatusFound)
	}),
})

mux := http.NewServeMux()
mux.Handle("/", auth)
mux.Handle("/private", secure(privateHandler))
http.ListenAndServe(":8080", mux)
```

Then, go to `http://localhost:8080/auth/login` to see the login page.

## Features

- Pluggable http.Handler that can be used with any go http server
- Customizable email sending mechanism
- Customizable email validation mechanism
- Customizable logo
- Customizable product name

### Roadmap

- Out-of-the-box support for generating time-bound tokens using TOTP (Time-Based One-Time Password).
- Customizable templates (Bring your own).
- Automatically handle token expiration based on time, providing security and convenience.
- Prevend CSRF attacks with token.
