package match

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/api/account"
	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/api/rating"
	db "github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/db/sqlc"
	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/game"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

var matchService *MatchService
var accService *account.AccountService
var ratService *rating.RatingService
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
	accService = account.NewService(testStore)
	ratService = rating.NewService(testStore)
	matchService = NewService(testStore, accService, ratService)
}

func TestMatchLifecycle(t *testing.T) {
	ctx := context.Background()

	p1, _ := accService.Create(ctx, &account.CreateAccountDTO{Mail: "p1@m.com", Username: "p1", Password: "password"})
	p2, _ := accService.Create(ctx, &account.CreateAccountDTO{Mail: "p2@m.com", Username: "p2", Password: "password"})

	id1 := *p1.AccountID
	id2 := *p2.AccountID

	summary := MatchSummaryDTO{
		IsNewMatch: true,
		P1ID:       id1,
		P2ID:       id2,
		Date:       time.Now(),
		GameInfo: &game.GameInfo{
			Winner:        "P1_WINS",
			Termination:   "LASER",
			MatchType:     "RANKED",
			BoardType:     1,
			TimeBase:      300000,
			TimeIncrement: 2,
			Log:           "Rf1%j1,j4,i4,i5,j5,j9%{300000};Tg6:f6%a8,a5,b5,b4,a4,a0%{300000};Rb4%j1,j4,i4,i5,j5,j9%{295000};Ri5xf6%a8,a5,b5,b4,e4,e5,f5,f6%{290000};",
		},
	}

	t.Run("FinalizeAndSaveMatch", func(t *testing.T) {
		rewards, err := matchService.FinalizeMatch(ctx, summary)
		assert.NoError(t, err)
		assert.NotNil(t, rewards)
		assert.Positive(t, rewards.P1XPDiff)
		assert.Positive(t, rewards.P1EloDiff)

		history, err := matchService.GetUserHistory(ctx, id1)
		assert.NoError(t, err)
		assert.NotEmpty(t, history)
	})

	t.Run("GetByID", func(t *testing.T) {
		matchID := int64(3)
		res, err := matchService.GetByID(ctx, matchID)
		if err == nil {
			assert.Equal(t, id1, res.P1ID)
		}
	})

	t.Run("PausedMatches", func(t *testing.T) {
		paused, err := matchService.GetPausedMatches(ctx, id1)
		assert.NoError(t, err)
		assert.Nil(t, paused)
	})

	t.Run("FinalizeMatch_InvalidUsers", func(t *testing.T) {
		summary.P1ID = 999999
		_, err := matchService.FinalizeMatch(ctx, summary)
		assert.Error(t, err)
	})
}

func TestMatchParsers(t *testing.T) {
	t.Run("parseMatches_Empty", func(t *testing.T) {
		res := parseMatches([]db.Match{})
		assert.Empty(t, res)
	})

	t.Run("parsePausedMatches_Empty", func(t *testing.T) {
		res := parsePausedMatches([]db.GetPausedMatchesRow{})
		assert.Empty(t, res)
	})

	t.Run("toCreateMatchParams", func(t *testing.T) {
		data := MatchSaveDTO{
			GameInfo: &game.GameInfo{BoardType: 1},
		}
		params := toCreateMatchParamsFromSaveDTO(data)
		assert.NotNil(t, params)
	})
}
