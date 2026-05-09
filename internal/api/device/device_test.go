package device

import (
	"context"
	"log"
	"os"
	"testing"

	db "github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/db/sqlc"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var deviceService *DeviceService
var testStore *db.Store

func init() {
	godotenv.Load("../../../.env")

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatalln("Error: DB_URL no encontrada")
	}

	ctx := context.Background()
	dbPool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalln("Error: No se pudo conectar con la base de datos")
	}
	if err := dbPool.Ping(ctx); err != nil {
		log.Fatalln("Error: No hay conexión real con la DB:", err)
	}

	testStore = db.NewStore(dbPool)
	deviceService = NewService(testStore)
}

func TestDeviceService_FullCoverage(t *testing.T) {
	ctx := context.Background()

	var userID int64 = 1
	testToken := "test-token-coverage"

	t.Run("RegisterDevice_Success", func(t *testing.T) {
		dto := DeviceDTO{Token: testToken}
		id, err := deviceService.RegisterDevice(ctx, dto, userID)
		assert.NoError(t, err)
		assert.NotZero(t, id)
	})

	t.Run("RegisterDevice_Error_Duplicate", func(t *testing.T) {
		dto := DeviceDTO{Token: testToken}
		_, err := deviceService.RegisterDevice(ctx, dto, userID)
		assert.Error(t, err)
	})

	t.Run("GetDevicesById_Success", func(t *testing.T) {
		tokens, err := deviceService.GetDevicesById(ctx, userID)
		assert.NoError(t, err)
		assert.Contains(t, tokens, testToken)
	})

	t.Run("GetDevicesById_Empty", func(t *testing.T) {
		tokens, err := deviceService.GetDevicesById(ctx, 999999) // ID inexistente
		assert.NoError(t, err)
		assert.Empty(t, tokens)
	})

	t.Run("DeleteDevice_Success", func(t *testing.T) {
		deleted, err := deviceService.DeleteDevice(ctx, testToken)
		assert.NoError(t, err)
		assert.Equal(t, testToken, deleted)
	})

	t.Run("DeleteDevice_Not_Found", func(t *testing.T) {
		deviceService.DeleteDevice(ctx, "non-existent-token")
		assert.True(t, true)
	})
}

func TestRegisterDuplicateToken(t *testing.T) {
	ctx := context.Background()
	token := "token-unico"
	dto := DeviceDTO{Token: token}

	_, err := deviceService.RegisterDevice(ctx, dto, 1)
	require.NoError(t, err)

	_, err = deviceService.RegisterDevice(ctx, dto, 1)
	assert.Error(t, err)

	deviceService.DeleteDevice(ctx, token)
}
