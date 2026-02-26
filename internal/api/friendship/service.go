package friendship

import (
	db "github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/db/sqlc"
)

type friendship struct{
	queries *db.Queries
}

func NewService(q *db.Queries) *friendshipService {
	return &friendshipService{queries: q}
}