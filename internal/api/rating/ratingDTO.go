package rating

import (
	db "github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/db/sqlc"
)

// DTOs para tratar con ratings

// Para mover ratings de usuario
type RatingDTO struct {
	UserID  int64      `json:"user_id" binding:"required"`
	EloType db.EloType `json:"elo_type" binding:"required"`
	Value   int32      `json:"value" binding:"required"`
}

type GenericRatingDto struct {
	EloType db.EloType `json:"elo_type" binding:"required"`
	Value   int32      `json:"value" binding:"required"`
}

// Para mover todos los ratings de golpe
type AllRatingsDTO struct {
	UserID   int64 `json:"user_id" binding:"required"`
	Blitz    int32 `json:"blitz" binding:"required"`
	Extended int32 `json:"extended" binding:"required"`
	Rapid    int32 `json:"rapid" binding:"required"`
	Classic  int32 `json:"classic" binding:"required"`
}

type RankUserDTO struct {
	UserID 		int64 `json:"user_id" binding:"required"`
	Username 	string `json:"username" binding:"required"`
	Level 		int32 `json:"level" binding:"required"`
	Avatar 		int32 `json:"avatar" binding:"required"`
	Rating		int32 `json:"rating" binding:"required"`
}

type GetRankingDTO struct {
	EloType string `json:"elo_type" binding:"required"`
}

type GetRankByIdDTO struct {
	EloType db.EloType `json:"elo_type" binding:"required"`
}

type RankingDTO struct {
	Rank int64 `json:"rank" binding:"required"`
}
 