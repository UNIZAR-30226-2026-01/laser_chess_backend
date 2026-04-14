package device

// Service que se encarga de la lógica de negocio relacionada con las cuentas
// de usuario

import (
	"context"

	db "github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/db/sqlc"
)

type DeviceService struct {
	store *db.Store
}

func NewService(s *db.Store) *DeviceService {
	return &DeviceService{store: s}
}

// Registra un nuevo dispositivo al usuario con id == DeviceID
func (s *DeviceService) RegisterDevice(ctx context.Context,
	token RegisterDeviceDTO, userID int64) (int64, error) {

	return s.store.RegisterDevice(ctx, db.RegisterDeviceParams{
		UserID: userID,
		Token:  token.Token,
	})

}

func (s *DeviceService) GetDevicesById(ctx context.Context,
	userID int64) ([]string, error) {

	return s.store.GetDevicesById(ctx, userID)

}

func (s *DeviceService) DeleteDevice(ctx context.Context,
	token string) (string, error) {
	return s.store.DeleteDevice(ctx, token)
}