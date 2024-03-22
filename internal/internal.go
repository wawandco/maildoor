package internal

import "embed"

//go:embed *.html
var templates embed.FS

//go:embed *.png
var Assets embed.FS
