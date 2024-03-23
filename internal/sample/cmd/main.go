package main

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/wawandco/maildoor"
	"github.com/wawandco/maildoor/internal/sample"
)


func main() {
	r := http.NewServeMux()

	// Auth handlers
	r.Handle("/auth/", maildoor.New(
		maildoor.UsePrefix("/auth/"),
		maildoor.AfterLogin(func(w http.ResponseWriter, r *http.Request) {
			expire := time.Now().AddDate(0, 0, 1)
    		cookie := http.Cookie{
      			Name: "sample",
         		Value: r.Context().Value("email").(string),
           		Path: "/",
             	Domain: r.Host,
                Expires: expire,
                RawExpires: expire.Format(time.UnixDate),
                MaxAge: 86400,
                HttpOnly: true,
                SameSite: http.SameSiteStrictMode,
      		}

            http.SetCookie(w, &cookie)
			http.Redirect(w, r, "/private", http.StatusFound)
		}),
	))

	secure := func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			_, err := r.Cookie("sample")
			if err != nil {
				http.Redirect(w, r, "/auth/login", http.StatusFound)
				return
			}

			next.ServeHTTP(w, r)
		}
	}

	// Private endpoint
	r.HandleFunc("/private", secure(sample.Private))
	r.HandleFunc("/", sample.Home)

	slog.Info("Server running on :3000")
	http.ListenAndServe(":3000", r)
}
