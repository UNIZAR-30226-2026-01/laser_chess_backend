package account

// Service que se encarga de la lógica de negocio relacionada con las cuentas
// de usuario

import (
	"context"

	db "github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/db/sqlc"
	"golang.org/x/crypto/bcrypt"
)

type AccountService struct {
	queries *db.Queries
}

func NewService(q *db.Queries) *AccountService {
	return &AccountService{queries: q}
}

// Crea una cuenta
// Primero hashea la contraseña
// Por ahora se inventa los items equipados por defecto,
// pero habrá que hacer que los ownee y los tenga equipados.
func (s *AccountService) CreateAccount(ctx context.Context, body CreateAccountDTO) (AccountDTO, error) {

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		return AccountDTO{}, err
	}

	res, err := s.queries.CreateAccount(ctx, db.CreateAccountParams{
		PasswordHash: string(passwordHash),
		Username:     body.Username,
		Mail:         body.Mail,

		// Por ahora forzamos que sean 1 y 2,
		// pero habrá que hacer algún tipo de consulta
		// o algo
		BoardSkin: 1,
		PieceSkin: 2,
	})

	if err != nil {
		return AccountDTO{}, err
	}

	//TODO: hacer que ownee los cosmeticos por defecto

	// Solo devuelve el AccountID
	return AccountDTO{AccountID: res.AccountID}, nil
}

func (s *AccountService) GetAccountByID(ctx context.Context, accountID int64) (AccountDTO, error) {
	res, err := s.queries.GetAccountByID(ctx, accountID)
	if err != nil {
		return AccountDTO{}, err
	}

	return AccountDTO{
		AccountID: res.AccountID,
		Mail:      &res.Mail,
		Username:  &res.Username,
		Level:     &res.Level,
		Xp:        &res.Xp,
		Money:     &res.Money,
		BoardSkin: &res.BoardSkin,
		PieceSkin: &res.PieceSkin,
	}, nil
}

func (s *AccountService) UpdateAccount(ctx context.Context, body AccountDTO) (AccountDTO, error) {
	res, err := s.queries.UpdateAccount(ctx, db.UpdateAccountParams{
		AccountID: body.AccountID,
		Username:  body.Username,
		BoardSkin: body.BoardSkin,
		PieceSkin: body.PieceSkin,
	})
	if err != nil {
		return AccountDTO{}, err
	}
	return AccountDTO{
		AccountID: res.AccountID,
		Mail:      &res.Mail,
		Username:  &res.Username,
		BoardSkin: &res.BoardSkin,
		PieceSkin: &res.PieceSkin,
	}, nil
}

func (s *AccountService) DeleteAccount(ctx context.Context, accountID int64) error {
	err := s.queries.DeleteAccount(ctx, accountID)
	if err != nil {
		return err
	}

	return nil
}
