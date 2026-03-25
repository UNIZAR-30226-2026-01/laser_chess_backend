-- name: CreateFriendship :exec
INSERT INTO friendship (
    user1_id, user2_id, accepted_1, accepted_2)
VALUES (
    $1, $2, $3, $4
);

-- name: GetFriendship :one
SELECT * FROM friendship 
WHERE ($1 = user1_id AND $2 = user2_id) OR ($1 = user2_id AND $2 = user1_id);

-- name: GetUserFriendships :many
SELECT friendship.user2_id AS user_id, account.username, account.level, account.avatar 
FROM friendship 
JOIN account ON friendship.user2_id = account.account_id
WHERE accepted_1 = TRUE AND accepted_2 = TRUE AND $1 = friendship.user1_id

UNION

SELECT friendship.user1_id AS user_id, account.username, account.level, account.avatar 
FROM friendship 
JOIN account ON friendship.user1_id = account.account_id
WHERE accepted_1 = TRUE AND accepted_2 = TRUE AND $1 = friendship.user2_id;


-- name: GetUserPendingSentFriendships :many
SELECT friendship.user2_id AS user_id, account.username, account.level, account.avatar 
FROM friendship 
JOIN account ON friendship.user2_id = account.account_id
WHERE accepted_1 = TRUE AND accepted_2 = FALSE AND $1 = friendship.user1_id

UNION

SELECT friendship.user1_id AS user_id, account.username, account.level, account.avatar 
FROM friendship 
JOIN account ON friendship.user1_id = account.account_id
WHERE accepted_1 = FALSE AND accepted_2 = TRUE AND $1 = friendship.user2_id;


-- name: GetUserPendingReceivedFriendships :many
SELECT friendship.user2_id AS user_id, account.username, account.level, account.avatar 
FROM friendship 
JOIN account ON friendship.user2_id = account.account_id
WHERE accepted_1 = FALSE AND accepted_2 = TRUE AND $1 = friendship.user1_id

UNION

SELECT friendship.user1_id AS user_id, account.username, account.level, account.avatar 
FROM friendship 
JOIN account ON friendship.user1_id = account.account_id
WHERE accepted_1 = TRUE AND accepted_2 = FALSE AND $1 = friendship.user2_id;

-- name: AcceptFriendship :exec
UPDATE friendship
SET 
    accepted_1 = CASE WHEN user1_id = $1 THEN TRUE ELSE accepted_1 END,
    accepted_2 = CASE WHEN user2_id = $1 THEN TRUE ELSE accepted_2 END
WHERE (user1_id = $1 AND user2_id = $2) OR (user1_id = $2 AND user2_id = $1);


-- name: DeleteFriendship :exec
DELETE FROM friendship 
WHERE ($1 = user1_id AND $2 = user2_id) OR ($1 = user2_id AND $2 = user1_id);

-- QUERIES ANTIGUAS, POR SI ACASO

-- --ERES USER 1
-- -- name: GetUser1PendingSentFriendship :many
-- SELECT user2_id FROM friendship
-- WHERE ($1 = user1_id AND accepted_1 = TRUE AND accepted_2 = FALSE);

-- --ERES USER 2
-- -- name: GetUser2PendingSentFriendship :many
-- SELECT user1_id FROM friendship
-- WHERE ($1 = user2_id AND accepted_1 = FALSE AND accepted_2 = TRUE);

-- --ERES USER 1
-- -- name: GetUser1PendingReceivedFriendship :many
-- SELECT user2_id FROM friendship
-- WHERE ($1 = user1_id AND accepted_1 = FALSE AND accepted_2 = TRUE);

-- --ERES USER 2
-- -- name: GetUser2PendingReceivedFriendship :many
-- SELECT user1_id FROM friendship
-- WHERE ($1 = user2_id AND accepted_1 = TRUE AND accepted_2 = FALSE);

---- name: SetFriendship :exec
--UPDATE friendship
--SET accepted_1 = TRUE, accepted_2 = TRUE
--WHERE (user1_id = LEAST($1::bigint, $2::bigint) AND user2_id = GREATEST($1::bigint, $2::bigint));
