package account

// Service que se encarga de la lógica de negocio relacionada con las cuentas
// de usuario

import (
	"context"

	db "github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/db/sqlc"
)

type AccountService struct {
	queries *db.Queries
}

func NewService(q *db.Queries) *AccountService {
	return &AccountService{queries: q}
}

func (s *AccountService) CreateAccount(ctx context.Context, body db.CreateAccountParams) (db.Account, error) {
	res, err := s.queries.CreateAccount(ctx, body)
	if err != nil {
		return db.Account{}, err
	}
	return res, nil
}

func (s *AccountService) GetAccountByID(ctx context.Context, accountID int64) (db.Account, error) {
	res, err := s.queries.GetAccountByID(ctx, accountID)
	if err != nil {
		return db.Account{}, err
	}
	return res, nil
}

func (s *AccountService) UpdateAccount(ctx context.Context, body db.UpdateAccountParams) (db.Account, error) {
	res, err := s.queries.UpdateAccount(ctx, body)
	if err != nil {
		return db.Account{}, err
	}
	return res, nil
}

func (s *AccountService) DeleteAccount(ctx context.Context, accountID int64) error {
	err := s.queries.DeleteAccount(ctx, accountID)
	if err != nil {
		return err
	}

	return nil
}
