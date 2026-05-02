package match

import (
	"time"

	db "github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/db/sqlc"
	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/game"
)

type PausedMatchDTO struct {
	MatchID         int64          `json:"match_id" binding:"required"`
	P1ID            int64          `json:"p1_id" binding:"required"`
	P2ID            int64          `json:"p2_id" binding:"required"`
	P1Username      string         `json:"p1_username"`
	P2Username      string         `json:"p2_username"`
	P1Elo           int32          `json:"p1_elo"`
	P2Elo           int32          `json:"p2_elo"`
	Date            time.Time      `json:"date"`
	Winner          db.Winner      `json:"winner"`
	Termination     db.Termination `json:"termination"`
	MatchType       db.MatchType   `json:"match_type"`
	Board           db.BoardType   `json:"board"`
	MovementHistory string         `json:"movement_history"`
	TimeBase        int32          `json:"time_base"`
	TimeIncrement   int32          `json:"time_increment"`
}

type MatchDTO struct {
	P1ID            int64          `json:"p1_id" binding:"required"`
	P2ID            int64          `json:"p2_id" binding:"required"`
	P1Elo           int32          `json:"p1_elo"`
	P2Elo           int32          `json:"p2_elo"`
	Date            time.Time      `json:"date"`
	Winner          db.Winner      `json:"winner"`
	Termination     db.Termination `json:"termination"`
	MatchType       db.MatchType   `json:"match_type"`
	Board           db.BoardType   `json:"board"`
	MovementHistory string         `json:"movement_history"`
	TimeBase        int32          `json:"time_base"`
	TimeIncrement   int32          `json:"time_increment"`
}

type MatchSaveDTO struct {
	IsNewMatch bool
	GameInfo   *game.GameInfo
	P1ID       int64
	P2ID       int64
	P1Elo      int32
	P2Elo      int32
	Date       time.Time
}

// Este DTO sirve para pasar toda la info desde la Room al Service
type MatchSummaryDTO struct {
	IsNewMatch bool
	GameInfo   *game.GameInfo
	P1ID       int64
	P2ID       int64
	Date       time.Time
}

// Para pasarle a la room las rewards y el elo
type MatchRewardsDTO struct {
	P1XPDiff    int32
	P2XPDiff    int32
	P1MoneyDiff int32
	P2MoneyDiff int32
	P1EloDiff   int32
	P2EloDiff   int32
}
