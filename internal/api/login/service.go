package login

// Service que se encarga de la lógica de negocio relacionada con las cuentas
// de usuario y los JWTs

import (
	"context"
	"time"

	account "github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/api/account"
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

// Coge el userID y contraseña hasheada de un user a partir de su mail
// o nombre de usuario
func (s *LoginService) getCredentials(ctx context.Context, body *LoginDTO) (int64, string, error) {
	if account.IsMail(body.Credential) {
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
		ExpiresAt: expiresAt,
	})
}

// Verifica las credenciales de un usuario,
// y genera access y refresh tokens
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

// Valida el refresh token, y si es correcto,
// crea una pareja nueva de refresh y access tokens
func (s *LoginService) Refresh(ctx context.Context, refreshToken string) (*LoginResult, error) {
	tokenHash := auth.HashToken(refreshToken)

	// Buscar token en bdd
	session, err := s.store.GetRefreshSession(ctx, tokenHash)
	if err != nil {
		return nil, apierror.ErrUnauthorized
	}

	// Si el refresh ha expirado no se hace nada
	// El user tendra que hacer login de nuevo
	if time.Now().After(session.ExpiresAt) {
		s.store.DeleteRefreshSession(ctx, tokenHash)
		return nil, apierror.ErrUnauthorized
	}

	// Volver a generar tokens
	newAccessToken, newRefreshToken, err := generateAccessRefreshToken(session.AccountID)
	if err != nil {
		return nil, err
	}

	// Borramos el antiguo y creamos nuevo
	s.store.DeleteRefreshSession(ctx, tokenHash)

	err = s.store.CreateRefreshSession(ctx, db.CreateRefreshSessionParams{
		AccountID: session.AccountID,
		TokenHash: auth.HashToken(newRefreshToken),
		ExpiresAt: time.Now().Add(auth.RefreshTokenTTL),
	})

	return &LoginResult{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil

}
