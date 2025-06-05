[![report card](https://goreportcard.com/badge/github.com/wawandco/maildoor)](https://goreportcard.com/report/github.com/wawandco/maildoor)

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
- Custom renderer functions for login and code entry pages

### Custom Renderers

Maildoor now supports custom renderer functions that allow you to completely customize the appearance of the login and code entry pages. You can provide your own HTML templates while still leveraging maildoor's authentication logic.

```go
auth := maildoor.New(
	maildoor.ProductName("My App"),
	maildoor.Logo("https://example.com/logo.png"),

	// Custom login page renderer
	maildoor.LoginRenderer(func(data maildoor.Attempt) (string, error) {
		// Return your custom HTML for the login page
		html := `<html>...your custom login page...</html>`
		return html, nil
	}),

	// Custom code entry page renderer
	maildoor.CodeRenderer(func(data maildoor.Attempt) (string, error) {
		// Return your custom HTML for the code entry page
		html := `<html>...your custom code page...</html>`
		return html, nil
	}),

	// Other options...
)
```

The `maildoor.Attempt` struct contains:
- `Logo` - URL of the logo image
- `Icon` - URL of the icon image
- `ProductName` - Name of your product
- `Email` - The email address (available in code renderer)
- `Error` - Error message if any validation failed
- `Code` - The verification code (context-dependent)

### Token Storage

Maildoor uses configurable token storage to manage authentication codes. By default, it uses in-memory storage, but you can provide custom implementations for Redis, databases, or other backends.

```go
// Use default in-memory storage (no expiration)
auth := maildoor.New(
	maildoor.ProductName("My App"),
	// ... other options
)

// Use in-memory storage with expiration
tokenStorage := maildoor.NewInMemoryTokenStorage(5 * time.Minute)
auth := maildoor.New(
	maildoor.WithTokenStorage(tokenStorage),
	maildoor.ProductName("My App"),
	// ... other options
)

// Use custom storage (implement TokenStorage interface)
customStorage := &MyRedisStorage{}
auth := maildoor.New(
	maildoor.WithTokenStorage(customStorage),
	maildoor.ProductName("My App"),
	// ... other options
)
```

### Roadmap

- Out of the box time bound token generation
- Time based token expiration out the box
- Prevend CSRF attacks with token
