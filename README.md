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
- Customizable logo and product name
- **Custom login page templates** - Complete control over login page appearance
- Responsive design with mobile-friendly defaults
- Built-in error handling and validation

## Custom Templates

Maildoor now supports fully customizable login page templates. You can replace the default layout and/or login form with your own HTML and CSS:

```go
// Custom login form template
customLogin := `{{define "yield"}}
<div class="my-custom-login">
    <h1>Welcome to {{.ProductName}}</h1>
    <form action="{{prefixedPath "/email"}}" method="POST">
        <input type="email" name="email" placeholder="Enter your email" required>
        <button type="submit">Sign In</button>
    </form>
    {{if ne .Error ""}}
        <div class="error">{{.Error}}</div>
    {{end}}
</div>
{{end}}`

// Custom layout template
customLayout := `<!DOCTYPE html>
<html>
<head>
    <title>{{block "title" .}}{{.ProductName}}{{end}}</title>
    <style>
        body { background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); }
        .container { max-width: 400px; margin: 50px auto; }
    </style>
</head>
<body>
    <div class="container">
        {{block "yield" .}}{{end}}
    </div>
</body>
</html>`

auth := maildoor.New(
    maildoor.CustomLayoutTemplate(customLayout),
    maildoor.CustomLoginTemplate(customLogin),
    maildoor.ProductName("MyApp"),
    // ... other options
)
```

**Available template data:**
- `.Logo` - Logo URL
- `.Icon` - Icon URL
- `.ProductName` - Product name
- `.Error` - Error message (if any)
- `.Email` - User email (in error states)

**Template functions:**
- `prefixedPath` - Generate URLs with correct prefix

See [CUSTOM_TEMPLATES.md](CUSTOM_TEMPLATES.md) for detailed documentation and examples.

### Roadmap

- Custom token storage mechanism
- Out of the box time bound token generation
- Time based token expiration out the box
- Prevent CSRF attacks with token
