package maildoor

import (
	_ "embed"
	"net/http"
)

// login function renders the login page
func (h handler) login(w http.ResponseWriter, r *http.Request) {
	err := buildTemplate("templates/login.html", w, struct {
		Action  string
		Logo    string
		Favicon string
	}{
		Action:  h.sendPath(),
		Logo:    h.product.LogoURL,
		Favicon: h.product.FaviconURL,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
