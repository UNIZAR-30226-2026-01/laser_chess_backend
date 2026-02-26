package account

// DTOs para tratar con accounts

// Para crear una cuenta
type CreateAccountDTO struct {
	Password  string `json:"password" binding:"required"`
	Mail      string `json:"mail" binding:"required"`
	Username  string `json:"username" binding:"required"`
	BoardSkin int32  `json:"board_skin" binding:"required"`
	PieceSkin int32  `json:"piece_skin" binding:"required"`
}

// Para mandar/recibir un user al/del frontend
// Solo es obligatorio el userID, el resto es opcional
type AccountDTO struct {
	AccountID int64   `json:"account_id" binding:"required"`
	Mail      *string `json:"mail"`
	Username  *string `json:"username"`
	Level     *int32  `json:"level"`
	Xp        *int32  `json:"xp"`
	Money     *int32  `json:"money"`
	BoardSkin *int32  `json:"board_skin"`
	PieceSkin *int32  `json:"piece_skin"`
}
