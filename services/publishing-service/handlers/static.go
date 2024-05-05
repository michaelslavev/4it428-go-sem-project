package handlers

import (
	"embed"
)

var (
	//go:embed static/post_template.html
	PostTemplate embed.FS
)
