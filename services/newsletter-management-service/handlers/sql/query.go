package sql

import (
	_ "embed"
)

var (
	//go:embed queries/ListNewsletters.sql
	ListNewslettersSQL string

	//go:embed queries/ListNewsletters.sql
	CreateNewsletterSQL string
)
