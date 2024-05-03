package model

import "time"

type Newsletter struct {
	ID          string    `db:"id"`
	CreatedAt   time.Time `db:"created_at"`
	Title       string    `db:"title"`
	Description string    `db:"description"`
	EditorID    string    `db:"editor_id"`
}

type NewNewsletter struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}
