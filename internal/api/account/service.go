package account

// Service que se encarga de la lógica de negocio relacionada con las cuentas
// de usuario

import (
	"context"
	"regexp"

	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/api/apierror"
	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/auth"
	db "github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/db/sqlc"
)

type AccountService struct {
	store *db.Store
}

func NewService(s *db.Store) *AccountService {
	return &AccountService{store: s}
}

// Crea una cuenta, haciendo las inicializaciones pertinentes:
//   - Crear tablas de rating
//   - Hacer que tenga los items por defecto
//
// Primero hashea la contraseña
// Por ahora se inventa los items equipados por defecto,
// pero habrá que hacer que los ownee y los tenga equipados.
func (s *AccountService) Create(ctx context.Context, body *CreateAccountDTO) (*AccountDTO, error) {

	if !IsMail(body.Mail) {
		return nil, apierror.ErrInvalidMailFormat
	}

	if len(body.Password) > 50 || len(body.Password) < 6 {
		return nil, apierror.ErrInvalidPasswordLenght
	}

	passwordHash, err := auth.HashPassword(body.Password)
	if err != nil {
		return nil, err
	}

	var res int64

	// Ejecutar en transaccion
	err = s.store.ExecTx(ctx, func(q *db.Queries) error {
		var errTx error

		res, errTx = q.CreateAccount(ctx, db.CreateAccountParams{
			PasswordHash: string(passwordHash),
			Username:     body.Username,
			Mail:         body.Mail,

			BoardSkin:    4,  // ID 4: BOARD_SKIN (Classic)
			PieceSkin:    1,  // ID 1: PIECE_SKIN (Classic)
			WinAnimation: 7,  // ID 7: WIN_ANIMATION (Classic)
			Avatar:       10, // ID 10: AVATAR (bot1_lila)
		})
		if errTx != nil {
			return errTx
		}

		// Inicializar ratings
		errTx = q.CreateRatings(ctx, res)
		if errTx != nil {
			return errTx
		}

		// Hacer que el usuario ownee los cosméticos por defecto
		defaultItemIDs := []int32{1, 4, 7, 10} // Los mismos IDs de arriba

		for _, itemID := range defaultItemIDs {
			errTx = q.CreateItemOwner(ctx, db.CreateItemOwnerParams{
				UserID: res,
				ItemID: itemID,
			})
			if errTx != nil {
				return errTx
			}
		}

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
		Mail:         body.Mail,
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

func (s *AccountService) GetStats(ctx context.Context, accountID int64) (*AccountStatsDTO, error) {
	stats, err := s.store.GetStats(ctx, accountID)

	return &AccountStatsDTO{
		Level: stats.Level,
		Xp:    stats.Xp,
		Money: stats.Money,
	}, err
}

func (s *AccountService) UpdateStats(
	ctx context.Context,
	accountID int64,
	body *AccountStatsDTO,
) error {
	err := s.store.UpdateStats(ctx, db.UpdateStatsParams{
		AccountID: accountID,
		Level:     body.Level,
		Money:     body.Money,
		Xp:        body.Xp,
	})

	return err
}

// Desactiva la cuenta del usuario con id == accountID
func (s *AccountService) Delete(ctx context.Context, accountID int64) error {
	return s.store.DeleteAccount(ctx, accountID)
}

// Cambia la contrasenha del usuario con id == accountID
func (s *AccountService) ChangePassword(ctx context.Context,
	dto ChangePasswordDTO, accountID int64) error {

	oldPasswordHash, err := s.store.GetPasswordById(ctx, accountID)
	if err != nil {
		return err
	}

	err = auth.VerifyPassword(oldPasswordHash, dto.OldPassword)
	if err != nil {
		return apierror.ErrUnauthorized
	}

	newPasswordHash, err := auth.HashPassword(dto.NewPassword)
	if err != nil {
		return err
	}

	_, err = s.store.ChangePassword(ctx, db.ChangePasswordParams{
		AccountID:    accountID,
		PasswordHash: newPasswordHash,
	})

	if err != nil {
		return err
	}

	return nil
}

// Comprueba si un string es una direccion de email o no
func IsMail(credential string) bool {
	emailRegex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	return emailRegex.MatchString(credential)
}
