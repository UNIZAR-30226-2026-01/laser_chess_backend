package account

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/api/apierror"
	db "github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/db/sqlc"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var accountService *AccountService
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
	accountService = NewService(testStore)
}

func TestIsMail(t *testing.T) {
	tests := []struct {
		email    string
		expected bool
	}{
		{"test@example.com", true},
		{"usuario.nombre@zaragoza.es", true},
		{"invalido-email", false},
		{"@sinusuario.com", false},
		{"con espacios@ejemplo.com", false},
	}

	for _, tt := range tests {
		t.Run(tt.email, func(t *testing.T) {
			assert.Equal(t, tt.expected, IsMail(tt.email))
		})
	}
}

func TestAccountLifecycle(t *testing.T) {
	ctx := context.Background()
	email := "tester@example.com"
	username := "tester_user"
	password := "password123"

	t.Run("CreateAccount Success", func(t *testing.T) {
		dto := &CreateAccountDTO{
			Mail:     email,
			Username: username,
			Password: password,
		}

		res, err := accountService.Create(ctx, dto)
		require.NoError(t, err)
		require.NotNil(t, res.AccountID)

		accountID := *res.AccountID

		t.Run("GetByID", func(t *testing.T) {
			acc, err := accountService.GetByID(ctx, accountID)
			assert.NoError(t, err)
			assert.Equal(t, username, *acc.Username)
			assert.Equal(t, email, *acc.Mail)
		})

		t.Run("GetIDByUsername", func(t *testing.T) {
			id, err := accountService.GetIDByUsername(ctx, username)
			assert.NoError(t, err)
			assert.Equal(t, accountID, id)
		})

		t.Run("GetUsernameByID", func(t *testing.T) {
			uname, err := accountService.GetUsernameByID(ctx, accountID)
			assert.NoError(t, err)
			assert.Equal(t, username, uname)
		})

		t.Run("UpdateAccount", func(t *testing.T) {
			newUsername := "updated_tester"
			updateDTO := &AccountDTO{
				Username: &newUsername,
			}
			updated, err := accountService.Update(ctx, accountID, updateDTO)
			assert.NoError(t, err)
			assert.Equal(t, newUsername, *updated.Username)
		})

		t.Run("Stats Operations", func(t *testing.T) {
			statsDTO := &AccountStatsDTO{
				Level: 10,
				Xp:    500,
				Money: 1000,
			}
			err := accountService.UpdateStats(ctx, accountID, statsDTO)
			assert.NoError(t, err)

			currStats, err := accountService.GetStats(ctx, accountID)
			assert.NoError(t, err)
			assert.Equal(t, int32(10), currStats.Level)
			assert.Equal(t, int32(1000), currStats.Money)
		})

		t.Run("ChangePassword", func(t *testing.T) {
			changeDTO := ChangePasswordDTO{
				OldPassword: password,
				NewPassword: "newpassword456",
			}
			err := accountService.ChangePassword(ctx, changeDTO, accountID)
			assert.NoError(t, err)

			wrongChangeDTO := ChangePasswordDTO{
				OldPassword: "wrongpassword",
				NewPassword: "evennewer123",
			}
			err = accountService.ChangePassword(ctx, wrongChangeDTO, accountID)
			assert.Equal(t, apierror.ErrUnauthorized, err)
		})

		t.Run("DeleteAccount", func(t *testing.T) {
			err := accountService.Delete(ctx, accountID)
			assert.NoError(t, err)
		})
	})
}

func TestCreateAccountValidation(t *testing.T) {
	t.Run("Invalid Email Format", func(t *testing.T) {
		dto := &CreateAccountDTO{Mail: "badmail", Password: "password123"}
		res, err := accountService.Create(context.Background(), dto)
		assert.Nil(t, res)
		assert.Equal(t, apierror.ErrInvalidMailFormat, err)
	})

	t.Run("Password Too Short", func(t *testing.T) {
		dto := &CreateAccountDTO{Mail: "valid@mail.com", Password: "123"}
		res, err := accountService.Create(context.Background(), dto)
		assert.Nil(t, res)
		assert.Equal(t, apierror.ErrInvalidPasswordLenght, err)
	})
}

func TestGetStatsEmpty(t *testing.T) {
	stats, err := accountService.GetStats(context.Background(), 0)
	assert.NoError(t, err)
	assert.Equal(t, int32(0), stats.Level)
}
