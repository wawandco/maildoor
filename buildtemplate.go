package maildoor

import (
	"embed"
	"fmt"
	"html/template"
	"io"
	"path/filepath"
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

	// Non HTML templates do not use layouts.
	if filepath.Ext(tpath) != ".html" {
		t, err := template.New(tpath).Parse(string(content))
		if err != nil {
			return err
		}

		return t.Execute(w, data)
	}

	layout, err := templates.ReadFile("templates/layout.html")
	if err != nil {
		return err
	}

	htmlTemplate, err := template.New("layout").Parse(string(layout))
	if err != nil {
		return fmt.Errorf("error parsing layout template: %w", err)
	}

	contents := fmt.Sprintf(`{{define "content"}}%s{{end}}`, string(content))
	htmlTemplate, err = template.Must(htmlTemplate.Clone()).Parse(contents)
	if err != nil {
		return fmt.Errorf("error parsing template %w", err)
	}

	err = htmlTemplate.Execute(w, data)
	if err != nil {
		return err
	}

	return nil
}
