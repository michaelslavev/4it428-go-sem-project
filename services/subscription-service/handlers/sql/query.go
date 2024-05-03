package sql

import (
	_ "embed"
)

var (
	//go:embed queries/ListSubscriptions.sql
	ListSubscriptionsSql string
)
