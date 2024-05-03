INSERT INTO
    posts (title, content, newsletter_id)
VALUES ($1, $2, $3)
RETURNING *;