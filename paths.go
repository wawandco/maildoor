package maildoor

import "path"

func (h handler) sendPath() string {
	return h.baseURL + path.Join(h.prefix, "/send/")
}

func (h handler) loginPath() string {
	return h.baseURL + path.Join(h.prefix, "/login/")
}

func (h handler) validatePath() string {
	return h.baseURL + path.Join(h.prefix, "/validate/")
}

func (h handler) stylesPath() string {
	return h.baseURL + path.Join(h.prefix, "/assets/styles/maildoor.min.css")
}
