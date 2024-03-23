package maildoor

// option for the auth
type option func(*auth)

func UsePrefix(prefix string) option {
	return func(a *auth) {
		a.prefix = prefix
	}
}
