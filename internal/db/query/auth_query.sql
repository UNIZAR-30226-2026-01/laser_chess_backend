-- name: CreateRefreshSession :exec
INSERT INTO "refresh_session" (account_id, token_hash, expires_at)
VALUES ($1, $2, $3);
