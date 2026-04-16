package item

import (
	"context"

	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/api/account"
	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/api/apierror"
	db "github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/db/sqlc"
)

type itemService struct {
	store          *db.Store
	accountService *account.AccountService
}

func NewService(s *db.Store, accounts *account.AccountService) *itemService {
	return &itemService{store: s, accountService: accounts}
}

/*
*
* Desc: Esta funcion llama a una query generada por sqlc que asigna un item a
una cuenta de usuario dado un JSON.
* --- Parametros ---
* ctx, context.Context - Es el contexto de gin.
* data, ItemOwnerDTO - Es el DTO del objeto a insertar.
* --- Resultados ---
* error - Es el error que se haya provocado en la consulta, o nil en caso
contrario.
*
*/
func (s *itemService) Create(
	ctx context.Context,
	accountID int64,
	itemID int32,
) error {

	err := s.store.ExecTx(ctx, func(q *db.Queries) error {
		// Cojemos la informacion del user
		accountInfo, errTx := s.accountService.GetByID(ctx, accountID)
		if errTx != nil {
			return errTx
		}

		// Cojemos la unformacion del item
		itemInfo, errTx := s.GetByID(ctx, itemID)
		if errTx != nil {
			return errTx
		}

		// Comprobamos que el user tenga suficiente dinero
		if *accountInfo.Money < itemInfo.Price {
			return apierror.ErrNotEnoughMoney
		}

		// Comprobamos que el user tenga el nivel suficiente
		if *accountInfo.Level < itemInfo.LevelRequisite {
			return apierror.ErrUserLevelTooLow
		}

		// Actualizamos el dinero del user
		errTx = s.accountService.UpdateStats(ctx, accountID,
			&account.AccountStatsDTO{
				Level: *accountInfo.Level,
				Xp:    *accountInfo.Xp,
				Money: *accountInfo.Money - itemInfo.Price,
			})

		// Creamos el objeto itemOwner
		errTx = s.store.CreateItemOwner(ctx, db.CreateItemOwnerParams{
			UserID: accountID,
			ItemID: itemID,
		})

		if errTx != nil {
			return errTx
		}

		return nil
	})

	return err
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
func (s *itemService) GetByID(ctx context.Context, itemID int32) (*ShopItemDTO, error) {

	res, err := s.store.GetShopItem(ctx, itemID)
	if err != nil {
		return nil, err
	}

	return &ShopItemDTO{
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

func (s *itemService) ListShopItems(ctx context.Context) ([]ShopItemDTO, error) {
	items, err := s.store.ListShopItems(ctx)
	return parseShopItemToDTO(items), err
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

func parseShopItemToDTO(data []db.ShopItem) []ShopItemDTO {
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
