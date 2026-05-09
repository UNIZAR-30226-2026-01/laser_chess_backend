package device

type DeviceDTO struct {
	Token string `json:"token" binding:"required"`
}
