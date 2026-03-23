package rt

// DTOs para partidas privadas

// Query params del endpoint de crear reto
// El id del challenger se saca del JWT
type CreateChallengeDTO struct {
	ChallengedUsername string `form:"username"       binding:"required"`
	Board              int    `form:"board" 		 binding:"required"`
	StartingTime       int    `form:"starting_time"  binding:"required"`
	TimeIncrement      int    `form:"time_increment" binding:"required"`
}

// Elemento de la lista de retos pendientes
type PendingChallengeDTO struct {
	ChallengerID       int64  `json:"challenger_id"`
	ChallengerUsername string `json:"challenger_username"`
	Board              int    `json:"board"`
	StartingTime       int    `json:"starting_time"`
	TimeIncrement      int    `json:"time_increment"`
}

// Query params del endpoint de aceptar reto
// El id del challenged se saca del JWT
type AcceptChallengeDTO struct {
	ChallengerUsername string `form:"username" binding:"required"`
}
