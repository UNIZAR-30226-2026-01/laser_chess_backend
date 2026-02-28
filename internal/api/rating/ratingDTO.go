package rating

// DTOs para tratar con ratings

// Para mover ratings de usuario
type RatingDTO struct {
	UserID  int64 `json:"user_id" bindig:"required"`
	EloType int32 `json:"elo_type" bindig:"required"`
	Value   int32 `json:"value" binding:"required"`
}

// Para mover todos los ratings de golpe
type AllRatingsDTO struct {
	BlitzElo RatingDTO
}
