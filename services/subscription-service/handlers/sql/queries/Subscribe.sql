INSERT INTO 
    subscriptions (newsletter_id, subscriber_id) 
VALUES 
    ($1, $2)
RETURNING *;