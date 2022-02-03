package maildoor

import (
	_ "embed"
	"net/http"
)

// ecodes holds the error messages for the supported error codes.
// these get rendered in the login page error box.
var ecodes = map[string]string{
	"E1": "Opps 😥  something happened while trying to find a user account with the given email. Please try again.",
	"E2": "We're sorry, the specified token has already expired. Please enter your email again to receive a new one.",
	"E3": "The token you have entered is invalid. Please enter your email again to receive a new one.",
	"E4": "🤔 Something was out of order with your previous login attempt. Please try again.",
}

// login function renders the login page, it also renders conditionally
// errors because when some of the other endpoints fail, it will redirect
// to this page.
func (h handler) login(w http.ResponseWriter, r *http.Request) {
	token, err := GenerateJWT(csrfDuration, []byte(h.csrfTokenSecret))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct {
		Action    string
		Logo      string
		Favicon   string
		CSRFToken string
		Error     string
	}{
		Action:    h.sendPath(),
		Logo:      h.product.LogoURL,
		Favicon:   h.product.FaviconURL,
		CSRFToken: token,
		Error:     ecodes[r.Form.Get("error")],
	}

	err = buildTemplate("templates/login.html", w, data)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
