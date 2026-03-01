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

/*
*
* Desc: Esta funcion llama a una query generada por sqlc que asigna un item a
una cuenta de usuario dado un JSON.
* --- Parametros ---
* ctx, context.Context - Es el contexto de gin.
* data, ItemOwnerDTO - Es el DTO del objeto a insertar.
* --- Resultados ---
* ItemOwnerDTO - Es un objeto de confirmacion que tan solo contiene los ids del
item y usuario del objeto creado.
* error - Es el error que se haya provocado en la consulta, o nil en caso
contrario.
*
*/
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

/*
*
* Desc: Esta funcion llama a una query generada por sqlc que busca un item dado
su id.
* --- Parametros ---
* ctx, context.Context - Es el contexto de gin.
* itemID, int32 - Es id del item a buscar.
* --- Resultados ---
* ShopItemDTO - Es la informacion del item a buscar.
* error - Es el error que se haya provocado en la consulta, o nil en caso
contrario.
*
*/
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

/*
*
* Desc: Esta funcion llama a una query generada por sqlc que busca los items de
una cuenta dado su id.
* --- Parametros ---
* ctx, context.Context - Es el contexto de gin.
* userID, int64 - Es id de la cuenta.
* --- Resultados ---
* []ShopItemDTO - Es la lista de items de la cuenta de usuario.
* error - Es el error que se haya provocado en la consulta, o nil en caso
contrario.
*
*/
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
