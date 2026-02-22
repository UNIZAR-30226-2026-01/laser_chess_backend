package match

import (
	"context"

	db "github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/db/sqlc"
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
		return db.Match{}, err
	}

	return res, nil
}

func (s *MatchService) GetUserHistory(ctx context.Context, userID int64) ([]db.Match, error) {
	res, err := s.queries.GetUserHistory(ctx, userID)
	if err != nil {
		return []db.Match{}, err
	}

	return res, nil
}
