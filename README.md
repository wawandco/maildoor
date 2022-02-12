# Maildoor

Maildoor is an email authentication library for Go (lang) it establishes a workflow to login users by emailing tokens to their email addresses instead of using a password. It provides an API to define application specific behaviors as part of the authentication process.

## Installation

This library is intended to be used as a dependency in your Go project. Installation implies go-getting the package with:

```sh
go get github.com/wawandco/maildoor@latest
```

And then using it accordingly in your app. See the Usage section for detailed instructions on usage.
## Usage

Maildoor instances satisfy the http.Handler interface and can be mounted into Mupliplexers. To initialize a Maildoor instance, use the New function:

```go
    auth, err := maildoor.New(maildoor.Options{
		CSRFTokenSecret: os.Getenv("CSRF_TOKEN_SECRET"),
		
		FinderFn:       finder,
        SenderFn:       sender,
		AfterLoginFn:   afterLogin,
		LogoutFn:       logout,
	})

	if err != nil {
		return nil, fmt.Errorf("error initializing maildoor: %w", err)
	}
```

After initializing the Maildoor instance, you can mount it into a multiplexer:

```go
    mux := http.NewServeMux()
    mux.Handle("/auth/", auth) // Set the prefix

    fmt.Println("Listening on port 8080")
	if err := http.ListenAndServe(":8080", server); err != nil {
		panic(fmt.Errorf("error starting server: %w", err))
	}
```


### Options

After seeing how to initialize the Maildoor Instance, lets dig a deeper into what some of these options mean.

#### CSRFTokenSecret

This option sets the secret used by the signin form to protect against CSRF attacks. We recommend to pull this value from an environment variable or secret storage.

#### FinderFn

The finder function is used to find a user by email address. The logic for looking up users is up to the application developer, but it should return an `Emailable` instance to be used on the signin flow.

The signature of the finder function is:

```go
func(string) (Emailable, error)
```

Where the string is the email address or token to identify the user.
#### SenderFn

Maildoor does not take care of sending your emails, instead it expects you to provide a function that will do this. This function will be called when a user requests a token to be sent to their email address and will be passed the message that needs to be send to the user.

The sender function signature is:

```go
func(*maildoor.Message) error
```

When this function returns an error the sign-in flow redirects the user to the login page with an error message.

#### AfterLoginFn

AfterLoginFn is a function that is called after a user has successfully logged in. It is passed the request instance, the response and user that has just logged in. Within this function typically the application does things like setting a session cookie and redirecting the user to a secure page. As with the sender function, its up to the application to decide what happens within the afterLogin function. 

Its signature is:

```go
func(w http.ResponseWriter, r *http.Request, user Emailable) error
```

#### LogoutFn

Similar than the afterLogin function, the logout function is called after a user has successfully logged out. It is passed the request instance, the response and user that has just logged out. Within this function typically the application does things like clearing the session cookie and redirecting the landing page. As with the afterLoginFn function, it's up to the application to decide what happens within the logout function.

Its signature is:

```go
func(w http.ResponseWriter, r *http.Request) error
```

#### Product

Product allows to set some product related settings for the signin flow. This helps branding the pages rendered to the user. The product can specify the name of the product, the logo and the favicon.

### The HTTP Endpoints

Maildoor is an http.Handler, which means it receives requests and responds to them. The Maildoor handler is mounted on a prefix, which is set by the application developer. Under that prefix the handler responds to the following endpoints:

#### GET:/auth/login
This is the login form. It renders a form with a CSRF token and a submit button. In here the user is asked to enter their email address.

#### POST:/auth/send
This endpoint is hit by the login form. It receives the email address and the CSRF token from the user and upon confirmation with the `FinderFn` it sends a link with token to the user's email address.
#### GET:/auth/validate
This endpoint is where the email link is validated. It receives the token from the URL and it validates the token. If the token is valid, it runs the `AfterLoginFn` function.
#### DELETE:/auth/logout
This endpoint is intended to be used by the app to logout the user. It just run the `LogoutFn` function.
    
### Sample Application

Within the sample folder you can find a go application that illustrates the usage of Maildoor. To run it from the command line you can use:

```sh
go run ./sample/cmd/sample
```
## FAQ

- I do not use SMTP for sending, what should I do?
- How to I customize the email logo and product?
- Can I change the email copy (Subject or content)?
- I don't want to use JWT for my tokens, what should I do?
- What should I do in the `AfterLoginFn` hook?
- How do I secure my application to prevent unauthorized access?

## Guiding Principles

- Use standard Go library as much as possible to avoid external dependencies.
- Application logic should live in the application, not in the library.
## TODO

- [x] Cover with tests
- [x] CSRF on the login form. 
- [x] Error messages
- [x] Custom Logger
- [x] Write Usage
- [x] Sample Go application
- [x] List and describe the http endpoints
- [ ] Answer FAQ
- [ ] Default logo and favicon
- [ ] SMTP senderFn
- [ ] Authentication Middleware ❓
- [ ] Error pages (500 and 404)


