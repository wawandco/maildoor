package main

import (
	"net/http"

	"github.com/wawandco/maildoor"
	"github.com/wawandco/maildoor/internal/sample"
)


func main() {
	r := http.NewServeMux()
	r.Handle("/auth/", maildoor.New())
	r.HandleFunc("/private", sample.Private)

	http.ListenAndServe(":3000", r)
}
