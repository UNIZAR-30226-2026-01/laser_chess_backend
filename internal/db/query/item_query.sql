-- name: CreateItemOwner :one
INSERT INTO item_owner (
    user_id, item_id
)
VALUES (
    $1, $2
)
RETURNING *;

-- name: GetUserItems :many
SELECT shop_item.item_id, price, level_requisite, item_type::ITEM_TYPE, is_default FROM item_owner 
JOIN shop_item ON shop_item.item_id = item_owner.item_id
WHERE user_id = $1;

-- name: GetShopItem :one
SELECT * FROM shop_item
WHERE item_id = $1 LIMIT 1;