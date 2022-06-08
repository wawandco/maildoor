package maildoor

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// send email if the user exists otherwise still say we have
// sent it, not to give an idea of existing/non-existing users.
func (h *handler) send(w http.ResponseWriter, r *http.Request) {
	valid, err := ValidateJWT(r.FormValue("CSRFToken"), []byte(h.csrfTokenSecret))
	if err != nil || !valid {
		http.Redirect(w, r, h.loginPath()+"?error=E4", http.StatusSeeOther)

		return
	}

	email := r.Form.Get("email")
	user, err := h.finderFn(email)
	if err != nil {
		http.Redirect(w, r, h.loginPath()+"?error=E1", http.StatusSeeOther)

		return
	}

	if user != nil {
		// only send the email if the user exists
		tt, err := h.tokenManager.Generate(user)
		if err != nil {
			http.Redirect(w, r, h.loginPath()+"?error=E6", http.StatusSeeOther)

			return
		}

		loginLink := fmt.Sprintf("%v?token=%v&email=%v", h.validatePath(), tt, url.QueryEscape(user.EmailAddress()))
		mm, err := h.composeMessage(user, loginLink)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = h.senderFn(mm)
		if err != nil {
			http.Redirect(w, r, h.loginPath()+"?error=E5", http.StatusSeeOther)

			return
		}
	}

	err = buildTemplate("templates/emailsent.html", w, struct {
		Title        string
		LoginPath    string
		EmailAddress string
		Favicon      string
		StylesPath   string
	}{
		Title:        "Authentication Email Sent",
		LoginPath:    h.loginPath(),
		EmailAddress: email,
		Favicon:      h.product.FaviconURL,

		StylesPath: h.stylesPath(),
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *handler) composeMessage(user Emailable, link string) (*Message, error) {
	mm := &Message{
		To:      user.EmailAddress(),
		Subject: "Your login link to " + h.product.Name,
	}

	data := struct {
		Product   string
		Logo      string
		Year      string
		LoginLink string
		BaseURL   string
	}{
		Product:   h.product.Name,
		Logo:      h.product.LogoURL,
		Year:      time.Now().Format("2006"),
		LoginLink: link,
		BaseURL:   h.baseURL,
	}

	bb := bytes.NewBuffer([]byte{})
	err := buildTemplate("templates/message.txt.email", bb, data)
	if err != nil {
		return nil, err
	}

	bplain := bb.Bytes()

	bb = bytes.NewBuffer([]byte{})
	err = buildTemplate("templates/message.html.email", bb, data)
	if err != nil {
		return nil, err
	}

	bhtml := bb.Bytes()
	mm.addBody("text/html", bhtml)
	mm.addBody("text/plain", bplain)

	return mm, nil
}
