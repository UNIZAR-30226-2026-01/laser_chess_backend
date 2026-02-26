-- Queries públicas desde endpoints

-- name: CreateAccount :one
INSERT INTO account (
    password_hash, mail, username, 
    board_skin, piece_skin
)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetAccountByID :one
SELECT * FROM account
WHERE account_id = $1 AND is_deleted = FALSE LIMIT 1;


-- solo cambia cosas qué se pueden cambiar por el user
-- usa coalesce para solo actualizar los params que no 
--  son null en la query
-- name: UpdateAccount :one
UPDATE account
SET 
    username = COALESCE(sqlc.narg('username'), username),
    board_skin = COALESCE(sqlc.narg('board_skin'), board_skin),
    piece_skin = COALESCE(sqlc.narg('piece_skin'), piece_skin)
WHERE account_id = $1 AND is_deleted = FALSE
RETURNING *;


-- Queries privadas que sólo se llaman desde dentro del sistema

-- name: GetAccountByMail :one
SELECT account_id, password_hash FROM account
WHERE mail = $1 AND is_deleted = FALSE;

-- name: GetAccountByUsername :one
SELECT account_id, password_hash FROM account
WHERE username = $1 AND is_deleted = FALSE;

-- name: DeleteAccount :exec
UPDATE account
SET
    mail = 'deleted_' || account_id::text,
    username = 'deleted_' || account_id::text,
    password_hash = '',
    is_deleted = TRUE
WHERE account_id = $1 AND is_deleted = FALSE;

-- name: UpdatePassword :exec
UPDATE account
SET
    password_hash = $2
WHERE account_id = $1 AND is_deleted = FALSE;

-- name: UpdateMail :one
UPDATE account
SET
    mail = $2
WHERE account_id = $1 AND is_deleted = FALSE
RETURNING *;

-- name: UpdateStats :one
UPDATE account
SET
    level = $2,
    xp = $3,
    money = $4
WHERE account_id = $1 AND is_deleted = FALSE
RETURNING *;
