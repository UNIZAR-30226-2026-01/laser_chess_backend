-- name: CreateRatings :many
INSERT INTO rating (
    user_id, elo_type, value
)
VALUES 
(
    $1, 'blitz',   $2
),
(
    $1, 'extended',  $3
),
(
    $1, 'rapid',   $4
),
(
    $1, 'classic', $5
)
RETURNING *;

-- name: GetAllElos :many
SELECT * FROM rating
WHERE user_id = $1;

-- name: GetBlitzElo :one
SELECT * FROM rating
WHERE user_id = $1 AND elo_type = 'blitz'::elo_type;

-- name: GetExtendedElo :one
SELECT * FROM rating
WHERE user_id = $1 AND elo_type = 'extended'::elo_type;

-- name: GetRapidElo :one
SELECT * FROM rating
WHERE user_id = $1 AND elo_type = 'rapid'::elo_type;

-- name: GetClassicElo :one
SELECT * FROM rating
WHERE user_id = $1 AND elo_type = 'classic'::elo_type;

-- name: UpdateRating :one
UPDATE rating
SET
    value = $3
WHERE user_id = $1 AND elo_type = $2
RETURNING *;