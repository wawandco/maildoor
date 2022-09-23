package maildoor

import (
	"net/http"
	"path"
	"strings"
)

// handler takes care of processing different actions against the maildoor
// server, such as login, send, validate, logout, most of these involve calling
// the corresponding functions provided by the host application.
type handler struct {
	prefix          string
	baseURL         string
	csrfTokenSecret string

	// Product settings
	product productConfig

	finderFn     func(token string) (Emailable, error)
	senderFn     func(message *Message) error
	afterLoginFn func(w http.ResponseWriter, r *http.Request, user Emailable) error
	logoutFn     func(w http.ResponseWriter, r *http.Request) error

	tokenManager TokenManager
	logger       Logger

	// Serves the static assets such as css and images
	assetsServer http.Handler

	valueEncoder valueEncoder
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.logger.Info(r.Method, ":", r.URL.Path)

	r.URL.Path = strings.TrimSuffix(r.URL.Path, "/")
	if !strings.HasPrefix(r.URL.Path, h.prefix) {
		r.URL.Path = path.Join(h.prefix, r.URL.Path)
	}

	err := r.ParseForm()
	if err != nil {
		h.logger.Errorf("error parsing form: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Overriding the request method to allow for browsers
	// to do DELETE/PATCH/PUT requests.
	if r.Form.Get("_method") != "" {
		h.logger.Infof("Request method upgraded to be %v", r.Form.Get("_method"))
		r.Method = r.Form.Get("_method")
	}

	if r.URL.Path == h.prefix {
		http.Redirect(w, r, path.Join(h.prefix, "/login/"), http.StatusFound)
		return
	}

	if r.URL.Path == path.Join(h.prefix, "/login/") && r.Method == http.MethodGet {
		h.login(w, r)

		return
	}

	if r.URL.Path == path.Join(h.prefix, "/send/") && r.Method == http.MethodPost {
		h.send(w, r)

		return
	}

	if r.URL.Path == path.Join(h.prefix, "/validate/") && r.Method == http.MethodGet {
		h.validate(w, r)

		return
	}

	if r.URL.Path == path.Join(h.prefix, "/logout/") && r.Method == http.MethodDelete {
		h.logout(w, r)

		return
	}

	if strings.HasPrefix(r.URL.Path, path.Join(h.prefix, "assets")) && r.Method == http.MethodGet {
		// Trimming the prefix to get the path of the asset
		r.URL.Path = strings.Replace(r.URL.Path, path.Join(h.prefix), "", 1)
		h.assetsServer.ServeHTTP(w, r)

		return
	}

	http.NotFound(w, r)
}

func (h handler) CookieValue(r *http.Request) (string, error) {
	v, err := r.Cookie(DefaultCookieName)
	if err != nil {
		return "", err
	}

	return h.valueEncoder.Decode(v.Value)
}
