package auth

// En este fichero se definen funciones para crear y validar JWTs

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/api/apierror"
	"github.com/golang-jwt/jwt/v5"
)

// Habrá que pillarlo del .env
var jwtSecret = []byte("estoesunsecretonolomiresporfa")

// Time To Live del access token en minutos
var AccessTokenTTL time.Duration = 15 * time.Minute
var RefreshTokenTTL time.Duration = 7 * 24 * time.Hour

// Genera un JWT de corta duracion
func GenerateAccessToken(accountID int64) (string, error) {

	expirationTime := time.Now().Add(AccessTokenTTL)

	claims := jwt.MapClaims{
		"sub": accountID,
		"exp": expirationTime.Unix(),
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// Crea un string aleatorio seguro de 64 caracteres
func GenerateRefreshToken() (string, error) {
	// rellena un array con bytes random
	bytes := make([]byte, 32)

	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}

// Hashea el refresh token
func HashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

// Valida el access token y devuelve el accountID
func ValidateAccessToken(tokenString string) (int64, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, apierror.ErrInvalidToken
		}
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return 0, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, apierror.ErrInvalidToken
	}

	sub, ok := claims["sub"].(float64)
	if !ok {
		return 0, apierror.ErrInvalidToken
	}

	return int64(sub), nil
}
