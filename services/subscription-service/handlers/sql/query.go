package sql

import (
	_ "embed"
)

var (
	//go:embed queries/ListSubscriptions.sql
	ListSubscriptionsSql string
	//go:embed queries/Subscribe.sql
	SubscribeSql string
	//go:embed queries/Unsubcribe.sql
	UnsubcribeSql string
)
