package docs

import (
	"embed"
)

//go:embed *.html
var SwaggerTemplate embed.FS
