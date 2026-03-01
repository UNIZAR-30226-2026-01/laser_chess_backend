-- name: CreateRefreshSession :exec
INSERT INTO "refresh_session" (account_id, token_hash, expires_at)
VALUES ($1, $2, $3);

-- name: GetRefreshSession :one
SELECT * FROM "refresh_session" WHERE token_hash = $1 LIMIT 1;

-- name: DeleteRefreshSession :exec
DELETE FROM "refresh_session" WHERE token_hash = $1;

-- name: CountSessionsByAccount :one
SELECT count(*) FROM "refresh_session" 
WHERE account_id = $1;

-- name: DeleteOldestSession :exec
DELETE FROM "refresh_session"
WHERE token_hash = (
    SELECT token_hash 
    FROM "refresh_session" r
    WHERE r.account_id = $1
    ORDER BY created_at ASC
    LIMIT 1
);
