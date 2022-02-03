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

	err = buildTemplate("templates/login.html", w, struct {
		Action    string
		Logo      string
		Favicon   string
		CSRFToken string
	}{
		Action:    h.sendPath(),
		Logo:      h.product.LogoURL,
		Favicon:   h.product.FaviconURL,
		CSRFToken: token,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
