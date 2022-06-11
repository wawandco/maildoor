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

func defaultAfterLogin(w http.ResponseWriter, r *http.Request, user Emailable) error {
	// Sets the DefaultCookieName cookie so the user can pass the
	// authenticated middleware.
	cookie := &http.Cookie{
		Name:     DefaultCookieName,
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Now().Add(defaultCookieDuration),
		Secure:   true,
		Value:    user.EmailAddress(),
	}

	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/private", http.StatusSeeOther)

	return nil
}

// defaultLogout clears the maildoor cookie and redirects to the root.
func defaultLogout(w http.ResponseWriter, r *http.Request) error {
	cookie := &http.Cookie{
		Name:     DefaultCookieName,
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Now().Add(-1),
		Secure:   true,

		// This will expire the cookie.
		MaxAge: -1,
	}

	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/", http.StatusSeeOther)

	return nil
}
