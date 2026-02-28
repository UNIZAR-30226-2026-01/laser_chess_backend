-- name: CreateRatings :many
INSERT INTO rating (
    user_id, elo_type, value
)
VALUES 
(
    $1, elo_type.blitz,   $2
),
(
    $1, elo_type.bullet,  $3
),
(
    $1, elo_type.rapid,   $4
),
(
    $1, elo_type.classic, $5
)
RETURNING *;

-- name: GetAllElos :many
SELECT * FROM rating
WHERE user_id = $1;

-- name: GetBlitzElo :one
SELECT * FROM rating
WHERE user_id = $1 AND elo_type = elo_type.blitz;

-- name: GetBulletElo :one
SELECT * FROM rating
WHERE user_id = $1 AND elo_type = elo_type.bullet;

-- name: GetRapidElo :one
SELECT * FROM rating
WHERE user_id = $1 AND elo_type = elo_type.rapid;

-- name: GetClassicElo :one
SELECT * FROM rating
WHERE user_id = $1 AND elo_type = elo_type.classic;