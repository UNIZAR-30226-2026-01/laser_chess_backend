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
-- name: UpdateAccount :one
UPDATE account
SET 
    username = $2,
    board_skin = $3,
    piece_skin = $4
WHERE account_id = $1 AND is_deleted = FALSE
RETURNING *;


-- Queries privadas que sólo se llaman desde dentro del sistema
-- TODO: cambiar contraseña, mail y actualizar cosas varias
