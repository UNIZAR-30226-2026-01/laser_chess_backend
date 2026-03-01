package login

// DTOs para el login

// Datos de entrada del login
type LoginDTO struct {
	Credential string `json:"credential" binding:"required"`
	Password   string `json:"password" binding:"required"`
}

// Respuesta del login (el refresh va por httpOnly cookie)
type LoginResponseDTO struct {
	AccessToken string `json:"access_token" binding:"required"`
}

// Tipo interno para el service
type LoginResult struct {
	AccessToken  string
	RefreshToken string
}
