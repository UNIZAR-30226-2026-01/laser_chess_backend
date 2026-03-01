package friendship

import (
	"context"

	db "github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/db/sqlc"
)

type friendshipService struct {
	store *db.Store
}

func NewService(s *db.Store) *friendshipService {
	return &friendshipService{store: s}
}

func (s *friendshipService) Create(
	ctx context.Context, data db.CreateFriendshipParams) (db.Friendship, error) {
	return s.store.CreateFriendship(ctx, data)
}

