package maildoor

import (
	"net/http"
)

// validate receives the token from the request as a parameter and
// calls the validator function to check if the token is valid. In case
// it is valid, it calls the afterLogin function to let the application
// control what happens after that.
func (h handler) validate(w http.ResponseWriter, r *http.Request) {
	token := r.Form.Get("token")
	email := r.Form.Get("email")

	valid, err := h.tokenManager.Validate(token)
	if err != nil || !valid {
		http.Redirect(w, r, h.loginPath()+"?error=E3", http.StatusSeeOther)

		return
	}

	user, err := h.finderFn(email)
	if err != nil {
		http.Redirect(w, r, h.loginPath()+"?error=E1", http.StatusSeeOther)
		return
	}

	if user == nil {
		http.Redirect(w, r, h.loginPath()+"?error=E7", http.StatusSeeOther)
		return
	}

	err = h.afterLoginFn(w, r, user)
	if err != nil {
		http.Redirect(w, r, h.loginPath()+"?error=E7", http.StatusSeeOther)
		return
	}
}
