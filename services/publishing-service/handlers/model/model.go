package model

import "time"

type Post struct {
	ID              string    `db:"id"`
	CreatedAt       time.Time `db:"created_at"`
	Title           string    `db:"title"`
	Content         string    `db:"content"`
	NewsletterID    string    `db:"newsletter_id"`
}