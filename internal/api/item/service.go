package item

import (
	"context"

	db "github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/db/sqlc"
)

type itemService struct {
	store *db.Store
}

func NewService(s *db.Store) *itemService {
	return &itemService{store: s}
}

func (s *itemService) Create(ctx context.Context, data ItemOwnerDTO) (ItemOwnerDTO, error) {

	res, err := s.store.CreateItemOwner(ctx, db.CreateItemOwnerParams{
		ItemID: data.ItemID,
		UserID: data.UserID,
	})

	if err != nil {
		return ItemOwnerDTO{}, err
	}

	return ItemOwnerDTO{UserID: res.UserID, ItemID: res.ItemID}, nil
}

func (s *itemService) GetByID(ctx context.Context, itemID int32) (ShopItemDTO, error) {

	res, err := s.store.GetShopItem(ctx, itemID)
	if err != nil {
		return ShopItemDTO{}, err
	}

	return ShopItemDTO{
		ItemID:         res.ItemID,
		Price:          res.Price,
		LevelRequisite: res.LevelRequisite,
		ItemType:       res.ItemType,
		IsDefault:      res.IsDefault,
	}, nil
}

func (s *itemService) GetUserItems(ctx context.Context, userID int64) ([]ShopItemDTO, error) {

	res, err := s.store.GetUserItems(ctx, userID)
	if err != nil {
		return []ShopItemDTO{}, err
	}

	return parseUserItems(res), nil
}

// Funcion auxiliar: pasar de db.GetUserItemsRow a ShopItemDTO
func parseUserItems(data []db.GetUserItemsRow) []ShopItemDTO {
	var res []ShopItemDTO

	print(len(data))

	for _, value := range data {
		res = append(res, ShopItemDTO{
			ItemID:         value.ItemID,
			Price:          value.Price,
			LevelRequisite: value.LevelRequisite,
			ItemType:       value.ItemType,
			IsDefault:      value.IsDefault,
		})
	}

	return res
}
