-- name: CreateRatings :exec
INSERT INTO rating (user_id, elo_type)
VALUES 
    ($1, 'blitz'),
    ($1, 'extended'),
    ($1, 'rapid'),
    ($1, 'classic');

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

-- name: UpdateRating :exec
UPDATE rating
SET value = $3,
    deviation = $4,
    volatility = $5,
    last_played_at = NOW()
WHERE user_id = $1 AND elo_type = $2;
