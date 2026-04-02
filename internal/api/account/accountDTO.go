package account

// DTOs para tratar con accounts

// Para crear una cuenta
type CreateAccountDTO struct {
	Password string `json:"password" binding:"required"`
	Mail     string `json:"mail" binding:"required"`
	Username string `json:"username" binding:"required"`
}

// Para mandar/recibir un user al/del frontend
// Solo es obligatorio el userID, el resto es opcional
type AccountDTO struct {
	AccountID    *int64  `json:"account_id"`
	Mail         *string `json:"mail,omitempty"`
	Username     *string `json:"username,omitempty"`
	Level        *int32  `json:"level,omitempty"`
	Xp           *int32  `json:"xp,omitempty"`
	Money        *int32  `json:"money,omitempty"`
	BoardSkin    *int32  `json:"board_skin,omitempty"`
	PieceSkin    *int32  `json:"piece_skin,omitempty"`
	WinAnimation *int32  `json:"win_animation,omitempty"`
	Avatar       *int32  `json:"avatar,omitempty"`
}

type RegisterDeviceDTO struct {
	token string `json:"token"`
}
