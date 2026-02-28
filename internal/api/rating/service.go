package rating

import (
	"context"

	db "github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/db/sqlc"
)

type RatingService struct {
	store *db.Store
}

func NewService(s *db.Store) *RatingService {
	return &RatingService{store: s}
}

func (s *RatingService) GetAllElosByID(ctx context.Context, userID int64) ([]db.Rating, error) {
	res, err := s.store.GetAllElos(ctx, userID)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *RatingService) GetBlitzEloByID(ctx context.Context, userID int64) (db.Rating, error) {
	res, err := s.store.GetBlitzElo(ctx, userID)
	if err != nil {
		return db.Rating{}, err
	}
	return res, nil
}

func (s *RatingService) GetBulletEloByID(ctx context.Context, userID int64) (db.Rating, error) {
	res, err := s.store.GetBulletElo(ctx, userID)
	if err != nil {
		return db.Rating{}, err
	}
	return res, nil
}

func (s *RatingService) GetRapidEloByID(ctx context.Context, userID int64) (db.Rating, error) {
	res, err := s.store.GetRapidElo(ctx, userID)
	if err != nil {
		return db.Rating{}, err
	}
	return res, nil
}

func (s *RatingService) GetClassicEloByID(ctx context.Context, userID int64) (db.Rating, error) {
	res, err := s.store.GetClassicElo(ctx, userID)
	if err != nil {
		return db.Rating{}, err
	}
	return res, nil
}
