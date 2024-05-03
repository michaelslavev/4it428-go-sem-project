package sql

import (
	_ "embed"
)

var (
	//go:embed queries/ListNewsletters.sql
	ListNewslettersSQL string

	//go:embed queries/CreateNewsletter.sql
	CreateNewsletterSQL string

	//go:embed queries/RenameNewsletter.sql
	RenameNewsletterSQL string

	//go:embed queries/DeleteNewsletter.sql
	DeleteNewsletterSQL string
)
