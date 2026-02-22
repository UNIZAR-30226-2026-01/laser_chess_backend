package match

import (
	"context"
	"errors"

	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/api/apierror"
	db "github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/db/sqlc"
	"github.com/jackc/pgx/v5"
)

type MatchService struct {
	queries *db.Queries
}

func NewService(q *db.Queries) *MatchService {
	return &MatchService{queries: q}
}

func (s *MatchService) Create(ctx context.Context, data db.CreateMatchParams) (db.Match, error) {
	return s.queries.CreateMatch(ctx, data)
}

func (s *MatchService) GetByID(ctx context.Context, matchID int64) (db.Match, error) {
	res, err := s.queries.GetMatch(ctx, matchID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.Match{}, apierror.ErrNotFound
		} else {
			return db.Match{}, err
		}
	}

	return res, nil
}

func (s *MatchService) GetUserHistory(ctx context.Context, userID int64) ([]db.Match, error) {
	res, err := s.queries.GetUserHistory(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []db.Match{}, apierror.ErrNotFound
		} else {
			return []db.Match{}, err
		}
	}

	return res, nil
}
