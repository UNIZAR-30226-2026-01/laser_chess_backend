-- name: CreatePlaceholder :one
INSERT INTO placeholder (data)
VALUES ($1)
ON CONFLICT (data) DO NOTHING
RETURNING *;

-- name: GetPlaceholder :one
SELECT * FROM placeholder
WHERE id = $1 LIMIT 1;
