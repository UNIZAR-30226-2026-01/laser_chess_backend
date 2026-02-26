package match

import (
	"context"

	db "github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/db/sqlc"
)

type MatchService struct {
	store *db.Store
}

func NewService(s *db.Store) *MatchService {
	return &MatchService{store: s}
}

func (s *MatchService) Create(ctx context.Context, data db.CreateMatchParams) (db.Match, error) {
	res, err := s.store.CreateMatch(ctx, data)
	if err != nil {
		return db.Match{}, err
	}
	return res, nil
}

func (s *MatchService) GetByID(ctx context.Context, matchID int64) (db.Match, error) {
	res, err := s.store.GetMatch(ctx, matchID)
	if err != nil {
		return db.Match{}, err
	}

	return res, nil
}

func (s *MatchService) GetUserHistory(ctx context.Context, userID int64) ([]db.Match, error) {
	res, err := s.store.GetUserHistory(ctx, userID)
	if err != nil {
		return nil, err
	}

	return res, nil
}
