package login

// Service que se encarga de la lógica de negocio relacionada con las cuentas
// de usuario

import (
	"context"
	"regexp"
	"time"

	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/api/apierror"
	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/auth"
	db "github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/db/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

type LoginService struct {
	store *db.Store
}

func NewService(s *db.Store) *LoginService {
	return &LoginService{store: s}
}

// Verifica las credenciales de un usuario,
func (s *LoginService) Login(ctx context.Context, body LoginDTO) (*LoginResult, error) {
	var res struct {
		accountID    int64
		passwordHash string
	}

	// ver si es mail o username
	// modularizar: GetByCredential o algo asi
	if isMail(body.Credential) {
		mailRes, err := s.store.GetAccountByMail(ctx, body.Credential)
		if err != nil {
			return nil, err
		}

		res.accountID = mailRes.AccountID
		res.passwordHash = mailRes.PasswordHash
	} else {
		usernameRes, err := s.store.GetAccountByUsername(ctx, body.Credential)
		if err != nil {
			return nil, err
		}

		res.accountID = usernameRes.AccountID
		res.passwordHash = usernameRes.PasswordHash
	}

	err := auth.VerifyPassword(res.passwordHash, body.Password)
	if err != nil {
		return nil, apierror.ErrUnauthorized
	}

	// Generar tokens
	accessToken, err := auth.GenerateAccessToken(res.accountID)
	if err != nil {
		return nil, err
	}

	refreshToken, err := auth.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	// Guardar refresh token en bdd
	expiresAt := time.Now().Add(7 * 24 * time.Hour)
	err = s.store.CreateRefreshSession(ctx, db.CreateRefreshSessionParams{
		AccountID: res.accountID,
		TokenHash: auth.HashToken(refreshToken),
		ExpiresAt: pgtype.Timestamptz{
			Time: expiresAt,
			Valid: true,
		}
	})
	if err != nil {
		return nil, err
	}

	return &LoginResult{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// Comprueba si un string es una direccion de email o no
func isMail(credential string) bool {
	emailRegex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	return emailRegex.MatchString(credential)
}
