-- name: CreateMatch :one
INSERT INTO match (
    p1_id, p2_id, p1_elo, p2_elo, date, winner, termination, match_type, board,
	movement_history, time_base, time_increment 
)
VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
)
RETURNING *;

-- name: GetMatch :one
SELECT * FROM match
WHERE match_id = $1 LIMIT 1;

-- name: GetUserHistory :many
SELECT * FROM match
WHERE p1_id = $1 OR p2_id = $1
ORDER BY date DESC;