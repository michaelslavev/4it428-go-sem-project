package sql

import (
	_ "embed"
)

var (
	//go:embed queries/ListPosts.sql
	ListPostsSql string
)
