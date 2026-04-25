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
    last_updated_at = NOW()
WHERE user_id = $1 AND elo_type = $2;

-- name: GetTopRankUsers :many
SELECT r.value, a.account_id, a.username, a.avatar
FROM rating r 
JOIN account a ON a.account_id = r.user_id
WHERE r.elo_type = $1
ORDER BY r.value DESC, r.user_id ASC LIMIT 100;

-- name: GetRankById :one
WITH user_score AS (
    SELECT r1.value 
    FROM rating r1
    WHERE r1.user_id = $2 AND r1.elo_type = $1
)
SELECT CAST(COUNT(*) + 1 AS BIGINT) AS rank
FROM rating r2
WHERE r2.elo_type = $1 
  AND (
      -- Contamos a los que tienen más puntos
      r2.value > (SELECT value FROM user_score)
      OR 
      -- O a los que tienen los mismos puntos pero menor ID (tu desempate)
      (r2.value = (SELECT value FROM user_score) AND r2.user_id < $2)
  );
