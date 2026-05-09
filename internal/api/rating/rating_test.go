package rating

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/api/account"
	db "github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/db/sqlc"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var ratingService *RatingService
var accountService *account.AccountService
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
	accountService = account.NewService(testStore)
	ratingService = NewService(testStore)
}

func TestRatingService(t *testing.T) {
	ctx := context.Background()

	resAcc, err := accountService.Create(ctx, &account.CreateAccountDTO{
		Mail:     "rating_tester@example.com",
		Username: "ratingtester",
		Password: "password123",
	})
	require.NoError(t, err)
	userID := *resAcc.AccountID

	t.Run("GetEloTypeFromBaseTime", func(t *testing.T) {
		assert.Equal(t, db.EloTypeBLITZ, GetEloTypeFromBaseTime(300000))
		assert.Equal(t, db.EloTypeRAPID, GetEloTypeFromBaseTime(900000))
		assert.Equal(t, db.EloTypeCLASSIC, GetEloTypeFromBaseTime(1800000))
		assert.Equal(t, db.EloTypeEXTENDED, GetEloTypeFromBaseTime(3600000))
	})

	t.Run("GetInitialElos", func(t *testing.T) {
		elos, err := ratingService.GetAllElosByID(ctx, userID)
		assert.NoError(t, err)
		assert.Equal(t, userID, elos.UserID)
		assert.Equal(t, int32(1500), elos.Blitz)
	})

	t.Run("GetEloByID_Variants", func(t *testing.T) {
		res, err := ratingService.GetEloByID(ctx, userID, 300000)
		assert.NoError(t, err)
		assert.Equal(t, db.EloTypeBLITZ, res.EloType)

		res, err = ratingService.GetEloByID(ctx, userID, 900000)
		assert.NoError(t, err)
		assert.Equal(t, db.EloTypeRAPID, res.EloType)

		res, err = ratingService.GetEloByID(ctx, userID, 1800000)
		assert.NoError(t, err)
		assert.Equal(t, db.EloTypeCLASSIC, res.EloType)

		res, err = ratingService.GetEloByID(ctx, userID, 3600000)
		assert.NoError(t, err)
		assert.Equal(t, db.EloTypeEXTENDED, res.EloType)
	})

	t.Run("GetEloByID_ZeroUser", func(t *testing.T) {
		res, err := ratingService.GetBlitzEloByID(ctx, 0)
		assert.NoError(t, err)
		assert.Equal(t, int32(1500), res.Value)

		res, err = ratingService.GetRapidEloByID(ctx, 0)
		assert.NoError(t, err)
		assert.Equal(t, int32(1500), res.Value)

		res, err = ratingService.GetClassicEloByID(ctx, 0)
		assert.NoError(t, err)
		assert.Equal(t, int32(1500), res.Value)

		res, err = ratingService.GetExtendedEloByID(ctx, 0)
		assert.NoError(t, err)
		assert.Equal(t, int32(1500), res.Value)
	})

	t.Run("UpdateElo", func(t *testing.T) {
		newRating := &RatingDTO{
			UserID:     userID,
			EloType:    db.EloTypeBLITZ,
			Value:      1600,
			Deviation:  50,
			Volatility: 0.06,
		}
		err := ratingService.UpdateEloByID(ctx, newRating)
		assert.NoError(t, err)

		updated, err := ratingService.GetBlitzEloByID(ctx, userID)
		assert.NoError(t, err)
		assert.Equal(t, int32(1600), updated.Value)
	})

	t.Run("RankingOperations", func(t *testing.T) {
		top, err := ratingService.GetTopRankUsers(ctx, "blitz")
		assert.NoError(t, err)
		assert.NotEmpty(t, top)

		rank, err := ratingService.GetRankById(ctx, "blitz", userID)
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, rank, int64(1))
	})
}

func TestAuxiliaryFunctions(t *testing.T) {
	t.Run("sqlcParamToDTO_Empty", func(t *testing.T) {
		res := sqlcParamToDTO([]db.Rating{})
		assert.Equal(t, &AllRatingsDTO{}, res)
	})

	t.Run("ParseRankingRow_Empty", func(t *testing.T) {
		res := ParseRankingRow([]db.GetTopRankUsersRow{})
		assert.Empty(t, res)
	})
}
