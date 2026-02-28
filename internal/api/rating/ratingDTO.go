package rating

import (
	db "github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/db/sqlc"
)

// DTOs para tratar con ratings

// Para mover ratings de usuario
type RatingDTO struct {
	UserID  int64      `json:"user_id" bindig:"required"`
	EloType db.EloType `json:"elo_type" bindig:"required"`
	Value   int32      `json:"value" binding:"required"`
}

type GenericRatingDto struct {
	EloType db.EloType `json:"elo_type" bindig:"required"`
	Value   int32      `json:"value" binding:"required"`
}

// Para mover todos los ratings de golpe
type AllRatingsDTO struct {
	UserID int64 `json:"user_id" bindig:"required"`
	Elo1   GenericRatingDto
	Elo2   GenericRatingDto
	Elo3   GenericRatingDto
	Elo4   GenericRatingDto
}
