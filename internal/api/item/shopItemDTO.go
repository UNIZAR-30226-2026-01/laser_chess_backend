package item

import (
	db "github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/db/sqlc"
)

// DTOs para tratar con accounts

// Para mandar/recibir un item_owner al/del frontend
// Son obligatorios ambos parametros
type ItemOwnerDTO struct {
	UserID int64 `json:"user_id" binding:"required"`
	ItemID int32 `json:"item_id" binding:"required"`
}

// Para mandar/recibir un shop_item al/del frontend
// El obligatorio el id del item
type ShopItemDTO struct {
	ItemID         int32       `json:"item_id" binding:"required"`
	Price          int32       `json:"price"`
	LevelRequisite int32       `json:"level_requisite" `
	ItemType       db.ItemType `json:"item_type"`
	IsDefault      bool        `json:"is_default"`
}
