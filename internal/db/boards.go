package db

import (
	db "github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/db/sqlc"
	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/game"
)

var IntToBoard = map[game.Board_T]db.BoardType{
	0: "ACE",
	1: "CURIOSITY",
	2: "SOPHIE",
	3: "GRAIL",
	4: "MERCURY",
}

var BoardToInt = map[db.BoardType]game.Board_T{
	"ACE":       0,
	"CURIOSITY": 1,
	"SOPHIE":    2,
	"GRAIL":     3,
	"MERCURY":   4,
}
