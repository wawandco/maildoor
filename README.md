# Maildoor

Maildoor is an email authentication library for Go (lang) it establishes a workflow to login users by emailing tokens to their email addresses instead of using a password. It provides an API to define application specific behaviors as part of the authentication process.

## Installation

This library is intended to be used as a dependency in your Go project. Installation implies go-getting the package with:

```sh
go get github.com/wawandco/maildoor@latest
```

And then using it accordingly in your app. See the Usage section for detailed instructions on usage.
## Usage

## FAQ

- I do not use SMTP for sending, what should I do?
- How to I customize the email logo and product?
- Can I change the email copy (Subject or content)?
- I don't want to use JWT for my tokens, what should I do?
- What should I do in the `AfterLoginFn` hook?
- How do I secure my application to prevent unauthorized access?

## Guiding Principles

- Use standard Go library as much as possible to avoid external dependencies.
## TODO

- [x] Cover with tests
- [x] CSRF on the login form. 
- [x] Error messages
- [x] Custom Logger
- [ ] Write Usage
- [ ] Sample Go application
- [ ] Answer FAQ
- [ ] Default logo and favicon
- [ ] SMTP senderFn
- [ ] Authentication Middleware ❓
- [ ] Error pages (500 and 404)


