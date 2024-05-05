package model

import "time"

type Subscription struct {
	ID              string    `db:"id"`
	CreatedAt       time.Time `db:"created_at"`
	NewsletterID    string    `db:"newsletter_id"`
	SubscriberID    string    `db:"subscriber_id"`
}

type Subscriber struct {
	Email string `db:"email"`
}

type Newsletter struct {
	Title string `db:"title"`
}