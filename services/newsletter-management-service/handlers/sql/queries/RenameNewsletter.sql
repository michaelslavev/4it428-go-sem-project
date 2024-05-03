UPDATE
    newsletters
SET
    title = $1
WHERE
    id = $2
    AND
    editor_id = $3
RETURNING *;