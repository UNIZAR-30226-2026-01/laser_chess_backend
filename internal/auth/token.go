package auth

// En este fichero se definen funciones para crear y validar JWTs

import (
	"crypto/rand"
	"encoding/hex"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = os.Getenv("JWT_SECRET")

// Time To Live del access token en minutos
var accessTokenTTL time.Duration = 15 * time.Minute

// Genera un JWT de corta duracion
func GenerateAccessToken(accountID int64) (string, error) {

	expirationTime := time.Now().Add(accessTokenTTL)

	claims := jwt.MapClaims{
		"sub": accountID,
		"exp": expirationTime,
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// Crea un string aleatorio seguro de 64 caracteres
func GenerateRefreshToken() (string, error) {
	// rellena un array con bytes random
	bytes := make([]byte, 32)
	rand.Read(bytes)

	return hex.EncodeToString(bytes), nil
}

// Valida el access token y devuelve el accountID
func ValidateAccessToken(tokenString string) (int64, error) {
	// TODO
}
