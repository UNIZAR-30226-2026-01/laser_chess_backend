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

-- name: GetPausedMatches :many
SELECT
    m.match_id,
    m.p1_id,
    m.p2_id,
    a1.username AS p1_username,
    a2.username AS p2_username,
    m.p1_elo,
    m.p2_elo,
    m.date,
    m.winner,
    m.termination,
    m.match_type,
    m.board,
    m.movement_history,
    m.time_base,
    m.time_increment
FROM match m
JOIN account a1 ON m.p1_id = a1.account_id
JOIN account a2 ON m.p2_id = a2.account_id
WHERE (m.p1_id = $1 OR m.p2_id = $1) AND m.termination = 'UNFINISHED'::termination
ORDER BY m.date DESC;

-- name: UpdateMatch :one
UPDATE match
SET p1_id = $1, p2_id = $2, p1_elo = $3, p2_elo = $4, date = $5, winner = $6, 
    termination = $7, match_type = $8, board = $9, movement_history = $10, 
    time_base = $11, time_increment = $12
WHERE match_id = $13
RETURNING *;
