# Maildoor

Maildoor is an email authentication library for Go (lang) it establishes a workflow to login users by emailing tokens to their email addresses instead of using a password. It provides an API to define application specific behaviors as part of the authentication process.

## Installation

This library is intended to be used as a dependency in your Go project. Installation implies go-getting the package with:

```sh
go get github.com/wawandco/maildoor@latest
```

And then using it acordingly in your app. See the Usage section for detailed instructions on usage.
## Example

## FAQ

## Guiding Principles

- Use standard Go library as much as possible to avoid external dependencies.
## TODO

- [x] Cover with tests
- [x] CSRF on the login form. 
- [x] Error messages
- [ ] Custom Logger
- [ ] Default logo and favicon
- [ ] Authentication Middleware
- [ ] SMTP senderFn
- [ ] Examples


