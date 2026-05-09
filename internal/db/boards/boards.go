package boards

import (
	db "github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/db/sqlc"
	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/game"
)

const BOARD_NUM = 5

var IntToBoard = map[game.Board_T]db.BoardType{
	0: "ACE",
	1: "CURIOSITY",
	2: "GRAIL",
	3: "MERCURY",
	4: "SOPHIE",
}

var BoardToInt = map[db.BoardType]game.Board_T{
	"ACE":       0,
	"CURIOSITY": 1,
	"GRAIL":     2,
	"MERCURY":   3,
	"SOPHIE":    4,
}
