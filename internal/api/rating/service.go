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

func (s *RatingService) GetAllElosByID(ctx context.Context, userID int64) (AllRatingsDTO, error) {
	res, err := s.store.GetAllElos(ctx, userID)
	if err != nil {
		return AllRatingsDTO{}, err
	}
	ratingsDTO := AllRatingsDTO{}
	ratingsDTO.UserID = res[0].UserID
	ratingsDTO.Elo1.EloType = res[0].EloType
	ratingsDTO.Elo1.Value = res[0].Value
	ratingsDTO.Elo2.EloType = res[1].EloType
	ratingsDTO.Elo2.Value = res[1].Value
	ratingsDTO.Elo3.EloType = res[2].EloType
	ratingsDTO.Elo3.Value = res[2].Value
	ratingsDTO.Elo4.EloType = res[3].EloType
	ratingsDTO.Elo4.Value = res[3].Value
	return ratingsDTO, nil
}

func (s *RatingService) GetBlitzEloByID(ctx context.Context, userID int64) (RatingDTO, error) {
	res, err := s.store.GetBlitzElo(ctx, userID)
	if err != nil {
		return RatingDTO{}, err
	}
	return RatingDTO{
		UserID:  res.UserID,
		EloType: res.EloType,
		Value:   res.Value,
	}, nil
}

func (s *RatingService) GetBulletEloByID(ctx context.Context, userID int64) (RatingDTO, error) {
	res, err := s.store.GetBulletElo(ctx, userID)
	if err != nil {
		return RatingDTO{}, err
	}
	return RatingDTO{
		UserID:  res.UserID,
		EloType: res.EloType,
		Value:   res.Value,
	}, nil
}

func (s *RatingService) GetRapidEloByID(ctx context.Context, userID int64) (RatingDTO, error) {
	res, err := s.store.GetRapidElo(ctx, userID)
	if err != nil {
		return RatingDTO{}, err
	}
	return RatingDTO{
		UserID:  res.UserID,
		EloType: res.EloType,
		Value:   res.Value,
	}, nil
}

func (s *RatingService) GetClassicEloByID(ctx context.Context, userID int64) (RatingDTO, error) {
	res, err := s.store.GetClassicElo(ctx, userID)
	if err != nil {
		return RatingDTO{}, err
	}
	return RatingDTO{
		UserID:  res.UserID,
		EloType: res.EloType,
		Value:   res.Value,
	}, nil
}
