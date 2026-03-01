-- name: CreateFriendship :one
INSERT INTO friendship (
    user1_id, user2_id, accepted_1, accepted_2)
VALUES (
    $1, $2, $3, $4
)
RETURNING *;

-- name: GetFriendship :one
SELECT * FROM friendship 
WHERE ($1 = user1_id AND $2 = user2_id) OR ($1 = user2_id AND $2 = user1_id);

-- name: GetUserFriendships :many
SELECT user2_id FROM friendship
WHERE accepted_1 = TRUE AND accepted_2 = TRUE AND $1 = user1_id;

--ERES USER 1
-- name: GetUser1PendingSentFriendship :many
SELECT user2_id FROM friendship
WHERE ($1 = user1_id AND accepted_1 = TRUE AND accepted_2 = FALSE);

--ERES USER 2
-- name: GetUser2PendingSentFriendship :many
SELECT user1_id FROM friendship
WHERE ($1 = user2_id AND accepted_1 = FALSE AND accepted_2 = TRUE);

--ERES USER 1
-- name: GetUser1PendingRecievedFriendship :many
SELECT user2_id FROM friendship
WHERE ($1 = user1_id AND accepted_1 = FALSE AND accepted_2 = TRUE);

--ERES USER 2
-- name: GetUser2PendingRecievedFriendship :many
SELECT user1_id FROM friendship
WHERE ($1 = user2_id AND accepted_1 = TRUE AND accepted_2 = FALSE);

-- name: SetFriendship :exec
UPDATE friendship
SET accepted_1 = TRUE, accepted_2 = TRUE
WHERE ($1 = user1_id AND $2 = user2_id) OR ($2 = user2_id AND $1 = user1_id);

-- name: DeleteFriendship :exec
DELETE FROM friendship 
WHERE ($1 = user1_id AND $2 = user2_id) OR ($1 = user2_id AND $2 = user1_id);

