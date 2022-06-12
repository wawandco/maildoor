package maildoor

import (
	"net/http"
	"time"
)

const (
	// DefaultCookieName is the of the cookie that
	// Maildoor will use to store the user's email address.
	// when the default login/logout functions are not overridden.
	DefaultCookieName     = "maildoor-cookie"
	defaultCookieDuration = 7 * 24 * time.Hour
)

func defaultAfterLogin(encoder valueEncoder) func(http.ResponseWriter, *http.Request, Emailable) error {
	return func(w http.ResponseWriter, r *http.Request, user Emailable) error {
		encoded, err := encoder.Encode(user.EmailAddress())
		if err != nil {
			return err
		}

		// Sets the DefaultCookieName cookie so the user can pass the
		// authenticated middleware.
		cookie := &http.Cookie{
			Name:    DefaultCookieName,
			Value:   encoded,
			Expires: time.Now().Add(defaultCookieDuration),

			Path:     "/",
			HttpOnly: true,
			Secure:   true,
		}

		http.SetCookie(w, cookie)
		http.Redirect(w, r, "/private", http.StatusSeeOther)

		return nil
	}
}

// defaultLogout clears the maildoor cookie and redirects to the root.
func defaultLogout(w http.ResponseWriter, r *http.Request) error {
	cookie := &http.Cookie{
		Name:     DefaultCookieName,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,

		// This will expire the cookie.
		MaxAge:  -1,
		Expires: time.Now().Add(-1),
	}

	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/", http.StatusSeeOther)

	return nil
}
