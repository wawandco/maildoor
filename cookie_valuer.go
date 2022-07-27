package maildoor

import "net/http"

type CookieValuer interface {
	CookieValue(r *http.Request) (string, error)
}
