package private

// DTOs para partidas privadas

// Query params del endpoint de crear reto
// El id del challenger se saca del JWT
type CreateChallengeDTO struct {
	ChallengedUsername *string `form:"username"`
	Board              *int    `form:"board"`
	StartingTime       *int32  `form:"starting_time"`
	TimeIncrement      *int32  `form:"time_increment"`
	MatchId            *int64  `form:"match_id"`
}

// Elemento de la lista de retos pendientes
type PendingChallengeDTO struct {
	ChallengerID       int64  `json:"challenger_id"`
	ChallengerUsername string `json:"challenger_username"`
	Board              int    `json:"board"`
	StartingTime       int32  `json:"starting_time"`
	TimeIncrement      int32  `json:"time_increment"`
}

// Query params del endpoint de aceptar reto
// El id del challenged se saca del JWT
type AcceptChallengeDTO struct {
	ChallengerUsername string `form:"username" binding:"required"`
}
