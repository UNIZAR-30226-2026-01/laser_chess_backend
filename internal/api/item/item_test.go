package item

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/api/account"
	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/api/apierror"
	db "github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/db/sqlc"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var itemService *ItemService
var accountService *account.AccountService
var testStore *db.Store

func init() {
	godotenv.Load("../../../.env")

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatalln("Error: DB_URL no encontrada")
	}

	ctx := context.Background()
	dbPool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalln("Error: No se pudo conectar con la base de datos")
	}

	testStore = db.NewStore(dbPool)
	accountService = account.NewService(testStore)
	itemService = NewService(testStore, accountService)
}

func TestShopItemService(t *testing.T) {
	ctx := context.Background()

	resAcc, err := accountService.Create(ctx, &account.CreateAccountDTO{
		Mail:     "shop_tester@example.com",
		Username: "shoptester",
		Password: "password123",
	})
	require.NoError(t, err)
	userID := *resAcc.AccountID

	t.Run("GetByID", func(t *testing.T) {
		item, err := itemService.GetByID(ctx, 2)
		assert.NoError(t, err)
		assert.Equal(t, int32(2), item.ItemID)
		assert.Equal(t, int32(500), item.Price)
		assert.Equal(t, int32(5), item.LevelRequisite)

		_, err = itemService.GetByID(ctx, 9999)
		assert.Error(t, err)
	})

	t.Run("ListShopItems", func(t *testing.T) {
		items, err := itemService.ListShopItems(ctx)
		assert.NoError(t, err)
		assert.NotEmpty(t, items)
	})

	t.Run("PurchaseFlow", func(t *testing.T) {
		t.Run("NotEnoughMoney", func(t *testing.T) {
			err := accountService.UpdateStats(ctx, userID, &account.AccountStatsDTO{
				Level: 20,
				Money: 0,
				Xp:    0,
			})
			require.NoError(t, err)

			err = itemService.Create(ctx, userID, 3)
			assert.Equal(t, apierror.ErrNotEnoughMoney, err)
		})

		t.Run("LevelTooLow", func(t *testing.T) {
			err := accountService.UpdateStats(ctx, userID, &account.AccountStatsDTO{
				Level: 1,
				Money: 5000,
				Xp:    0,
			})
			require.NoError(t, err)

			err = itemService.Create(ctx, userID, 3)
			assert.Equal(t, apierror.ErrUserLevelTooLow, err)
		})

		t.Run("SuccessfulPurchase", func(t *testing.T) {
			err := accountService.UpdateStats(ctx, userID, &account.AccountStatsDTO{
				Level: 15,
				Money: 2000,
				Xp:    0,
			})
			require.NoError(t, err)

			err = itemService.Create(ctx, userID, 21)
			assert.NoError(t, err)

			items, err := itemService.GetUserItems(ctx, userID)
			assert.NoError(t, err)
			found := false
			for _, it := range items {
				if it.ItemID == 21 {
					found = true
					break
				}
			}
			assert.True(t, found)
		})
	})

	t.Run("GetUserItems_Empty", func(t *testing.T) {
		items, err := itemService.GetUserItems(ctx, 999999)
		assert.NoError(t, err)
		assert.Empty(t, items)
	})
}

func TestParsers(t *testing.T) {
	t.Run("parseUserItems", func(t *testing.T) {
		data := []db.GetUserItemsRow{
			{ItemID: 1, Price: 0, LevelRequisite: 0, ItemType: "AVATAR", IsDefault: true},
		}
		res := parseUserItems(data)
		assert.Len(t, res, 1)
		assert.Equal(t, int32(1), res[0].ItemID)

		assert.Empty(t, parseUserItems([]db.GetUserItemsRow{}))
	})

	t.Run("parseShopItemToDTO", func(t *testing.T) {
		data := []db.ShopItem{
			{ItemID: 1, Price: 0, LevelRequisite: 0, ItemType: "AVATAR", IsDefault: true},
		}
		res := parseShopItemToDTO(data)
		assert.Len(t, res, 1)
		assert.Equal(t, int32(1), res[0].ItemID)

		assert.Empty(t, parseShopItemToDTO([]db.ShopItem{}))
	})
}
