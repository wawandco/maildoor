package sample

import "net/http"

func Private(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("🔒 Welcome to the private section"))
}
