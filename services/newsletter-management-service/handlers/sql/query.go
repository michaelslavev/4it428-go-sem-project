package sql

import (
	_ "embed"
)

var (
	//go:embed queries/ListNewsletters.sql
	ListNewsletters string
)
