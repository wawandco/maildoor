package internal

import (
	"embed"
	"path"
)

//go:embed *.html
var templates embed.FS

//go:embed *.png
var Assets embed.FS

func prefixedHelper(prefix string) func(string) string {
	return func(p string) string {
		return path.Join(prefix, p)
	}
}
