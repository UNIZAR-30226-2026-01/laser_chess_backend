package match

import (
	db "github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/db/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

type MatchDTO struct {
	P1ID            int64              `json:"p1_id" binding:"required"`
	P2ID            int64              `json:"p2_id" binding:"required"`
	P1Elo           int32              `json:"p1_elo"`
	P2Elo           int32              `json:"p2_elo"`
	Date            pgtype.Timestamptz `json:"date"`
	Winner          db.Winner          `json:"winner"`
	Termination     db.Termination     `json:"termination"`
	MatchType       db.MatchType       `json:"match_type"`
	Board           db.BoardType       `json:"board"`
	MovementHistory string             `json:"movement_history"`
	TimeBase        int32              `json:"time_base"`
	TimeIncrement   int32              `json:"time_increment"`
}
