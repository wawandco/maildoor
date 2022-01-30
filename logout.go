package maildoor

import "net/http"

func (h handler) logout(w http.ResponseWriter, r *http.Request) {
	err := h.logoutFn(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
