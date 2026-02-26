package login

// DTOs para el login

type LoginDTO struct {
	Credential string `json:"credential" binding:"required"`
	Password   string `json:"password" binding:"required"`
}
