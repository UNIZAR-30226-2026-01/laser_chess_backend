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

// Comprueba si un string es una direccion de email o no
func isMail(credential string) bool {
	emailRegex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	return emailRegex.MatchString(credential)
}

func (s *LoginService) getCredentials(ctx context.Context, body *LoginDTO) (int64, string, error) {
	if isMail(body.Credential) {
		mailRes, err := s.store.GetAccountByMail(ctx, body.Credential)
		if err != nil {
			return -1, "", err
		}

		return mailRes.AccountID, mailRes.PasswordHash, nil

	} else {
		usernameRes, err := s.store.GetAccountByUsername(ctx, body.Credential)
		if err != nil {
			return -1, "", err
		}

		return usernameRes.AccountID, usernameRes.PasswordHash, nil
	}
}

// Genera una pareja de access y refresh tokens
func generateAccessRefreshToken(accountID int64) (string, string, error) {

	accessToken, err := auth.GenerateAccessToken(accountID)
	if err != nil {
		return "", "", err
	}
	refreshToken, err := auth.GenerateRefreshToken()
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// Método auxiliar para generar y guardar el refresh token respetando el límite de 3
func (s *LoginService) saveNewSession(ctx context.Context, accountID int64, refreshToken string) error {
	count, err := s.store.CountSessionsByAccount(ctx, accountID)
	if err != nil {
		return err
	}

	// FIFO
	if count >= 3 {
		err = s.store.DeleteOldestSession(ctx, accountID)
		if err != nil {
			return err
		}
	}

	// Crear nueva sesion
	expiresAt := time.Now().Add(auth.RefreshTokenTTL)
	return s.store.CreateRefreshSession(ctx, db.CreateRefreshSessionParams{
		AccountID: accountID,
		TokenHash: auth.HashToken(refreshToken),
		ExpiresAt: pgtype.Timestamptz{Time: expiresAt, Valid: true},
	})
}

// Verifica las credenciales de un usuario,
func (s *LoginService) Login(ctx context.Context, body *LoginDTO) (*LoginResult, error) {

	// ver si es mail o username
	// modularizar: GetByCredential o algo asi
	accountID, passwordHash, err := s.getCredentials(ctx, body)
	if err != nil {
		return nil, apierror.ErrUnauthorized
	}

	err = auth.VerifyPassword(passwordHash, body.Password)
	if err != nil {
		return nil, apierror.ErrUnauthorized
	}

	// Generar tokens
	accessToken, refreshToken, err := generateAccessRefreshToken(accountID)
	if err != nil {
		return nil, err
	}

	// Guardar refresh token en bdd
	err = s.saveNewSession(ctx, accountID, refreshToken)
	if err != nil {
		return nil, err
	}

	return &LoginResult{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *LoginService) Refresh(ctx context.Context, refreshToken string) (*LoginResult, error) {
	tokenHash := auth.HashToken(refreshToken)

	// Buscar token en bdd
	session, err := s.store.GetRefreshSession(ctx, tokenHash)
	if err != nil {
		return nil, apierror.ErrUnauthorized
	}

	// Mirar si ha expirado
	if time.Now().After(session.ExpiresAt.Time) {
		s.store.DeleteRefreshSession(ctx, tokenHash)
		return nil, apierror.ErrUnauthorized
	}

	// Volver a generar tokens
	newAccessToken, newRefreshToken, err := generateAccessRefreshToken(session.AccountID)
	if err != nil {
		return nil, err
	}

	// Borramos el expirado y creamos nuevo
	s.store.DeleteRefreshSession(ctx, tokenHash)

	err = s.store.CreateRefreshSession(ctx, db.CreateRefreshSessionParams{
		AccountID: session.AccountID,
		TokenHash: auth.HashToken(newRefreshToken),
		ExpiresAt: pgtype.Timestamptz{Time: time.Now().Add(auth.RefreshTokenTTL), Valid: true},
	})

	return &LoginResult{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil

}
