package sql

import (
	_ "embed"
)

var (
	//go:embed queries/CheckSubscriptions.sql
	CheckSubscriptionsSql string
	//go:embed queries/Subscribe.sql
	SubscribeSql string
	//go:embed queries/Unsubcribe.sql
	UnsubcribeSql string
)
