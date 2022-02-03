package maildoor

import (
	_ "embed"
	"net/http"
	"time"
)

// login function renders the login page
func (h handler) login(w http.ResponseWriter, r *http.Request) {
	token, err := GenerateJWT(3*time.Minute, []byte(h.csrfTokenSecret))
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
	}

	ecode := r.Form.Get("error")
	switch ecode {
	case "E1":
		data.Error = "Opps 😥  something happened while trying to find a user account with the given email. Please try again."
	case "E2":
		data.Error = "We're sorry, the specified token has already expired. Please enter your email again to receive a new one."
	case "E3":
		data.Error = "The token you have entered is invalid. Please enter your email again to receive a new one."
	}

	err = buildTemplate("templates/login.html", w, data)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
