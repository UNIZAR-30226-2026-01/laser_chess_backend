package account

// Service que se encarga de la lógica de negocio relacionada con las cuentas
// de usuario

import (
	"context"

	db "github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/db/sqlc"
	"golang.org/x/crypto/bcrypt"
)

type AccountService struct {
	store *db.Store
}

func NewService(s *db.Store) *AccountService {
	return &AccountService{store: s}
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

	var res db.Account

	// Ejecutar en transaccion
	// Ahora no tiene sentido, pero lo tendra cuando hagamos lo de los items
	err = s.store.ExecTx(ctx, func(q *db.Queries) error {
		var errTx error

		res, errTx = q.CreateAccount(ctx, db.CreateAccountParams{
			PasswordHash: string(passwordHash),
			Username:     body.Username,
			Mail:         body.Mail,

			// Por ahora forzamos que sean 1 y 2,
			// pero habrá que hacer algún tipo de consulta
			// o algo
			BoardSkin: 1,
			PieceSkin: 2,
		})

		if errTx != nil {
			return errTx
		}

		//TODO: hacer que ownee los cosmeticos por defecto

		return nil
	})

	if err != nil {
		return AccountDTO{}, err
	}

	// Solo devuelve el AccountID
	return AccountDTO{AccountID: res.AccountID}, nil
}

func (s *AccountService) GetAccountByID(ctx context.Context, accountID int64) (AccountDTO, error) {
	res, err := s.store.GetAccountByID(ctx, accountID)
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
	res, err := s.store.UpdateAccount(ctx, db.UpdateAccountParams{
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
	err := s.store.DeleteAccount(ctx, accountID)
	if err != nil {
		return err
	}

	return nil
}
