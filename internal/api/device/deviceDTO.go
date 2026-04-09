package device

type RegisterDeviceDTO struct {
	Token string `json:"token" binding:"required"`
}