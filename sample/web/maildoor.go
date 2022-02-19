package web

import (
	"net/http"
	"time"

	"github.com/wawandco/maildoor"
)

func afterLogin(w http.ResponseWriter, r *http.Request, user maildoor.Emailable) error {
	// Sets the sample cookie so the user can pass the
	// authenticated middleware.
	cookie := &http.Cookie{
		Name:     "sample",
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
		Secure:   true,
		Value:    user.EmailAddress(),
	}

	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/private", http.StatusSeeOther)

	return nil
}

// logout function clears the sample cookie.
func logout(w http.ResponseWriter, r *http.Request) error {
	cookie := &http.Cookie{
		Name:     "sample",
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
