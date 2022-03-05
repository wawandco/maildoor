package maildoor

import (
	"net/http"
)

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
		Title     string
		Action    string
		Logo      string
		Favicon   string
		CSRFToken string
		Error     string

		StylesPath string
	}{
		Title:     "Login Page",
		Action:    h.sendPath(),
		Logo:      h.product.LogoURL,
		Favicon:   h.product.FaviconURL,
		CSRFToken: token,
		Error:     ecodes[r.Form.Get("error")],

		StylesPath: h.stylesPath(),
	}

	err = buildTemplate("templates/login.html", w, data)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
