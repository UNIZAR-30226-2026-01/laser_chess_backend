package account

// Service que se encarga de la lógica de negocio relacionada con las cuentas
// de usuario

import (
	"context"

	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/auth"
	db "github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/db/sqlc"
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
func (s *AccountService) Create(ctx context.Context, body *CreateAccountDTO) (*AccountDTO, error) {

	passwordHash, err := auth.HashPassword(body.Password)
	if err != nil {
		return nil, err
	}

	var res int64

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
			BoardSkin:    1,
			PieceSkin:    2,
			WinAnimation: 3,
			Avatar:       1,
		})

		if errTx != nil {
			return errTx
		}

		//TODO: hacer que ownee los cosmeticos por defecto

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Solo devuelve el AccountID
	return &AccountDTO{AccountID: &res}, nil
}

// Devuelve toda la info de la cuenta del user con id == accountID
func (s *AccountService) GetByID(ctx context.Context, accountID int64) (*AccountDTO, error) {
	res, err := s.store.GetAccountByID(ctx, accountID)
	if err != nil {
		return nil, err
	}

	return &AccountDTO{
		AccountID:    &res.AccountID,
		Mail:         &res.Mail,
		Username:     &res.Username,
		Level:        &res.Level,
		Xp:           &res.Xp,
		Money:        &res.Money,
		BoardSkin:    &res.BoardSkin,
		PieceSkin:    &res.PieceSkin,
		WinAnimation: &res.WinAnimation,
		Avatar:       &res.Avatar,
	}, nil
}

// Devuelve el username de la cuenta del user con id == accountID
func (s *AccountService) GetIDByUsername(ctx context.Context, username string) (int64, error) {
	res, err := s.store.GetAccountIDByUsername(ctx, username)
	if err != nil {
		return 0, err
	}

	return res, nil
}

// Devuelve el id de la cuenta del user con el username dado
func (s *AccountService) GetUsernameByID(ctx context.Context, ID int64) (string, error) {
	res, err := s.store.GetUsernameByID(ctx, ID)
	if err != nil {
		return "", err
	}

	return res, nil
}

// Actualiza el username o cosmeticos del usuario con
// id == accountID. Solo actualiza los campos no nulos
func (s *AccountService) Update(
	ctx context.Context,
	accountID int64,
	body *AccountDTO,
) (*AccountDTO, error) {
	res, err := s.store.UpdateAccount(ctx, db.UpdateAccountParams{
		AccountID:    accountID,
		Username:     body.Username,
		BoardSkin:    body.BoardSkin,
		PieceSkin:    body.PieceSkin,
		WinAnimation: body.WinAnimation,
		Avatar:       body.Avatar,
	})
	if err != nil {
		return nil, err
	}
	return &AccountDTO{
		AccountID:    &res.AccountID,
		Mail:         &res.Mail,
		Username:     &res.Username,
		Level:        &res.Level,
		Xp:           &res.Xp,
		Money:        &res.Money,
		BoardSkin:    &res.BoardSkin,
		PieceSkin:    &res.PieceSkin,
		WinAnimation: &res.WinAnimation,
		Avatar:       &res.Avatar,
	}, nil
}

// Desactiva la cuenta del usuario con id == accountID
func (s *AccountService) Delete(ctx context.Context, accountID int64) error {
	return s.store.DeleteAccount(ctx, accountID)
}

// Registra un nuevo dispositivo al usuario con id == accountID
func (s *AccountService) RegisterDevice(ctx context.Context,
	token RegisterDeviceDTO, accountID int64) (int64, error) {

	return s.store.RegisterDevice(ctx, db.RegisterDeviceParams{
		UserID: accountID,
		Token:  token.Token,
	})

}
