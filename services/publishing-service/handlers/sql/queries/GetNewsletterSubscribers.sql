SELECT
    u.id,
    u.email
FROM
    auth.users AS u
        JOIN
            subscriptions AS s ON u.id = s.subscriber_id
        JOIN
            newsletters AS n ON s.newsletter_id = n.id
WHERE
    n.id = $1
    AND
    n.editor_id = $2;
