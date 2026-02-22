package item

import (
	"context"
	"errors"

	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/api/apierror"
	db "github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/db/sqlc"
	"github.com/jackc/pgx/v5"
)

type itemService struct {
	queries *db.Queries
}

func NewService(q *db.Queries) *itemService {
	return &itemService{queries: q}
}

func (s *itemService) Create(ctx context.Context, data db.CreateItemOwnerParams) (db.ItemOwner, error) {
	return s.queries.CreateItemOwner(ctx, data)
}

func (s *itemService) GetByID(ctx context.Context, itemID int32) (db.ShopItem, error) {
	res, err := s.queries.GetShopItem(ctx, itemID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.ShopItem{}, apierror.ErrNotFound
		} else {
			return db.ShopItem{}, err
		}
	}

	return res, nil
}

func (s *itemService) GetUserItems(ctx context.Context, userID int64) ([]db.GetUserItemsRow, error) {
	res, err := s.queries.GetUserItems(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []db.GetUserItemsRow{}, apierror.ErrNotFound
		} else {
			return []db.GetUserItemsRow{}, err
		}
	}

	return res, nil
}
