-- name: CreateGreeter :one
INSERT INTO greeters (id, hello)
VALUES ($1, $2)
RETURNING *;

-- name: GetGreeter :one
SELECT * FROM greeters
WHERE id = $1 LIMIT 1;

-- name: ListGreeter :many
SELECT * FROM greeters
ORDER BY created_at DESC;