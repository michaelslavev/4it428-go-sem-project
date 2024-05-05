package handlers

import (
	"embed"
)

var (
	//go:embed static/mail_template.html
	SubscribedMailTemplate embed.FS
	//go:embed static/unsubscribed.html
	Unsubscribed embed.FS
)
