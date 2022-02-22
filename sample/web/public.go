package web

import "net/http"

func public(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	content, err := templates.ReadFile("templates/public.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	if _, err = w.Write(content); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}
}
