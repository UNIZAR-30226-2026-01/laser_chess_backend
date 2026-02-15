package placeholder

import (
	"context"
	"errors"

	db "github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/db/sqlc"
	"github.com/jackc/pgx/v5"
)

// Errores (habrá que poner todos los del proyecto en un paquete de errores)
var ErrNotFound = errors.New("resource not found")

type PlaceholderService struct {
	queries *db.Queries
}

func NewService(q *db.Queries) *PlaceholderService {
	return &PlaceholderService{queries: q}
}

func (s *PlaceholderService) Create(ctx context.Context, data string) (db.Placeholder, error) {
	return s.queries.CreatePlaceholder(ctx, data)
}

func (s *PlaceholderService) GetByID(ctx context.Context, id int32) (db.Placeholder, error) {
	res, err := s.queries.GetPlaceholder(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.Placeholder{}, ErrNotFound
		} else {
			return db.Placeholder{}, err
		}
	}

	return res, nil
}
