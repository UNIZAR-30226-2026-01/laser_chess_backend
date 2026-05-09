package friendship

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/api/account"
	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/api/apierror"
	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/api/device"
	db "github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/db/sqlc"
	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/sse"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

var friendshipService *FriendshipService
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

	testStore = db.NewStore(dbPool)

	deviceService := device.NewService(testStore)
	accService := account.NewService(testStore)

	fcm, err := sse.InitFirebase(deviceService)
	if err != nil {
		log.Println("Aviso: Firebase no inicializado, algunos eventos podrían fallar")
	}
	eventSystem := sse.InitSSE(fcm)

	friendshipService = NewService(testStore, eventSystem, accService)
}

func TestFriendshipLifecycle(t *testing.T) {
	ctx := context.Background()

	id1 := int64(1)
	id2 := int64(2)

	t.Run("Create Friendship Request", func(t *testing.T) {
		dto := &FriendshipDTO{
			SenderID:   &id1,
			ReceiverID: id2,
		}

		_ = friendshipService.DeleteFriendship(ctx, dto)

		err := friendshipService.Create(ctx, dto)
		assert.NoError(t, err)

		err = friendshipService.Create(ctx, dto)
		assert.Equal(t, apierror.ErrAlreadyFriends, err)
	})

	t.Run("Get Status and Pending", func(t *testing.T) {
		dto := &FriendshipDTO{SenderID: &id1, ReceiverID: id2}

		status, err := friendshipService.GetFriendshipStatus(ctx, dto)
		assert.NoError(t, err)
		assert.True(t, status.SenderAccept)
		assert.False(t, status.ReceiverAccept)

		sent, err := friendshipService.GetUserPendingSentFriendships(ctx, id1)
		assert.NoError(t, err)
		assert.NotEmpty(t, sent)

		received, err := friendshipService.GetUserPendingReceivedFriendships(ctx, id2)
		assert.NoError(t, err)
		assert.NotEmpty(t, received)

		count, err := friendshipService.GetUserPendingReceivedFriendshipsCount(ctx, id2)
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, count.Count, int32(1))
	})

	t.Run("Accept Friendship", func(t *testing.T) {
		dto := &FriendshipDTO{SenderID: &id2, ReceiverID: id1}

		err := friendshipService.AcceptFriendship(ctx, dto)
		assert.NoError(t, err)

		friends, err := friendshipService.GetUserFriendships(ctx, id1)
		assert.NoError(t, err)
		found := false
		for _, f := range friends {
			if f.UserID == id2 {
				found = true
				break
			}
		}
		assert.True(t, found)
	})

	t.Run("Delete Friendship", func(t *testing.T) {
		dto := &FriendshipDTO{SenderID: &id1, ReceiverID: id2}
		err := friendshipService.DeleteFriendship(ctx, dto)
		assert.NoError(t, err)

		friends, _ := friendshipService.GetUserFriendships(ctx, id1)
		assert.Empty(t, friends)
	})
}

func TestFriendshipErrors(t *testing.T) {
	ctx := context.Background()
	id := int64(1)

	t.Run("Sender Equals Receiver", func(t *testing.T) {
		dto := &FriendshipDTO{SenderID: &id, ReceiverID: id}
		err := friendshipService.Create(ctx, dto)
		assert.Equal(t, apierror.ErrBadRequest, err)

		_, err = friendshipService.GetFriendshipStatus(ctx, dto)
		assert.Equal(t, apierror.ErrBadRequest, err)
	})

	t.Run("Get Non Existent Friendship", func(t *testing.T) {
		dto := &FriendshipDTO{SenderID: &id, ReceiverID: 999999}
		_, err := friendshipService.GetFriendshipStatus(ctx, dto)
		assert.Error(t, err)
	})
}

func TestParsers(t *testing.T) {
	t.Run("ParseFriendshipRow Empty", func(t *testing.T) {
		res := ParseFriendshipRow([]db.GetUserFriendshipsRow{})
		assert.Nil(t, res)
	})

	t.Run("ParsePendingReceivedRow Empty", func(t *testing.T) {
		res := ParsePendingReceivedRow([]db.GetUserPendingReceivedFriendshipsRow{})
		assert.Nil(t, res)
	})
}
