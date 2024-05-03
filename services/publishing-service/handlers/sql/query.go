package sql

import (
	_ "embed"
)

var (
	//go:embed queries/ListPosts.sql
	ListPostsSql string
	//go:embed queries/CreatePost.sql
	CreatePostSql string
	//go:embed queries/GetNewsletterSubscribers.sql
	GetNewsletterSubscribersSql string
	//go:embed queries/GetNewsletter.sql
	GetNewsletter string
)
