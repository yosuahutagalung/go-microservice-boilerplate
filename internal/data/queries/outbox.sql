-- name: GetPendingOutboxEvents :many
SELECT id, topic, payload
FROM outbox_events
WHERE status = 'PENDING'
ORDER BY created_at ASC
LIMIT $1
FOR UPDATE SKIP LOCKED;

-- name: MarkOutboxEventPublished :exec
UPDATE outbox_events 
SET status = 'PUBLISHED' 
WHERE id = $1;

-- name: InsertOutboxEvent :exec
INSERT INTO outbox_events (topic, payload, status)
VALUES ($1, $2, 'PENDING');