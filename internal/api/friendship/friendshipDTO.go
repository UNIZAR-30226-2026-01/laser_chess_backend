package friendship

// DTOs para tratar con accounts

// Para crear una cuenta
type FreindshipDTO struct {
	AccountID int64 `json:"account_id" binding:"required"`
}


