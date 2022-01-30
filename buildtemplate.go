package maildoor

import (
	"embed"
	"html/template"
	"io"
)

var (
	//go:embed templates
	templates embed.FS
)

// buildTemplate for the passed template and data on a passed writer
// this is helpful to be able to render the templates in a generic way
// across different handlers.
func buildTemplate(tpath string, w io.Writer, data interface{}) error {
	content, err := templates.ReadFile(tpath)
	if err != nil {
		return err
	}

	t, err := template.New(tpath).Parse(string(content))
	if err != nil {
		return err
	}

	err = t.Execute(w, data)
	if err != nil {
		return err
	}

	return nil
}
