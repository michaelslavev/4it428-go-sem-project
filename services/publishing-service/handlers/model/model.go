package model

import "time"

type Post struct {
	ID              string    `db:"id"`
	CreatedAt       time.Time `db:"created_at"`
	Title           string    `db:"title"`
	Content         string    `db:"content"`
	NewsletterID    string    `db:"newsletter_id"`
}

type NewPost struct {
	Title        string `json:"title"`
	Content      string `json:"content"`
	NewsletterID string `json:"newsletterId"`
}

type Subscriber struct {
	Email string `db:"email"`
}

type Newsletter struct {
	ID          string    `db:"id"`
	CreatedAt   time.Time `db:"created_at"`
	Title       string    `db:"title"`
	Description string    `db:"description"`
	EditorID    string    `db:"editor_id"`
}