-- name: CreateRatings :exec
INSERT INTO rating (user_id, elo_type)
VALUES 
    ($1, 'BLITZ'),
    ($1, 'EXTENDED'),
    ($1, 'RAPID'),
    ($1, 'CLASSIC');

-- name: GetAllElos :many
SELECT * FROM rating
WHERE user_id = $1;

-- name: GetBlitzElo :one
SELECT * FROM rating
WHERE user_id = $1 AND elo_type = 'BLITZ'::elo_type;

-- name: GetExtendedElo :one
SELECT * FROM rating
WHERE user_id = $1 AND elo_type = 'EXTENDED'::elo_type;

-- name: GetRapidElo :one
SELECT * FROM rating
WHERE user_id = $1 AND elo_type = 'RAPID'::elo_type;

-- name: GetClassicElo :one
SELECT * FROM rating
WHERE user_id = $1 AND elo_type = 'CLASSIC'::elo_type;

-- name: UpdateRating :exec
UPDATE rating
SET value = $3,
    deviation = $4,
    volatility = $5,
    last_played_at = NOW()
WHERE user_id = $1 AND elo_type = $2;

-- name: GetTopRankUsers :many
SELECT r.value, a.account_id, a.username, a.level, a.avatar
FROM rating r 
JOIN account a ON a.account_id = r.user_id
WHERE r.elo_type = $1
ORDER BY r.value DESC LIMIT 100;

-- name: GetRankById :one
SELECT rank 
FROM (
    SELECT user_id, 
    ROW_NUMBER() OVER (ORDER BY value DESC) as rank
    FROM rating
    WHERE elo_type = $1 
) as rankings
WHERE user_id = $2;
