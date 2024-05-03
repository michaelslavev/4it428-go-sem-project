DELETE FROM 
    subscriptions
WHERE 
    newsletter_id = $1
AND
    subscriber_id = $2