package login

// Service que se encarga de la lógica de negocio relacionada con las cuentas
// de usuario

import (
	"context"
	"regexp"

	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/api/apierror"
	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/auth"
	db "github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/db/sqlc"
)

type LoginService struct {
	store *db.Store
}

func NewService(s *db.Store) *LoginService {
	return &LoginService{store: s}
}

func (s *LoginService) Login(ctx context.Context, body LoginDTO) error {
	var res struct {
		accountID    int64
		passwordHash string
	}

	// ver si es mail o username
	if isMail(body.Credential) {
		mailRes, err := s.store.GetAccountByMail(ctx, body.Credential)
		if err != nil {
			return err
		}

		res.accountID = mailRes.AccountID
		res.passwordHash = mailRes.PasswordHash
	} else {
		usernameRes, err := s.store.GetAccountByUsername(ctx, body.Credential)
		if err != nil {
			return err
		}

		res.accountID = usernameRes.AccountID
		res.passwordHash = usernameRes.PasswordHash
	}

	err := auth.VerifyPassword(res.passwordHash, body.Password)
	if err != nil {
		return apierror.ErrUnauthorized
	}

	return nil
}

func isMail(credential string) bool {
	emailRegex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	return emailRegex.MatchString(credential)
}
