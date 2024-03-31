package main

import (
	"log/slog"
	"net/http"

	"github.com/wawandco/maildoor/internal/sample"
)

func main() {
	r := http.NewServeMux()

	// Auth handlers
	r.Handle("/auth/", sample.Auth)

	// Application handlers
	r.HandleFunc("/private", sample.Private)
	r.HandleFunc("/{$}", sample.Home)

	slog.Info("Server running on :3000")
	http.ListenAndServe(":3000", r)
}
