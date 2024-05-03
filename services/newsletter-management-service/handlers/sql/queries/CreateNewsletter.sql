INSERT INTO
    newsletters (title, description, editor_id)
VALUES ($1, $2, $3)
RETURNING *;