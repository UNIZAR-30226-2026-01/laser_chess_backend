package login

// DTOs para el login

type LoginDTO struct {
	credential string `json:"credential" binding:"required"`
	password   string `json:"password" binding:"required"`
}
