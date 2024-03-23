package sample

import "net/http"

func Home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Public section"))
}
