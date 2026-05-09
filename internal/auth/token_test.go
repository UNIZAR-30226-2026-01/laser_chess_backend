package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateAndValidateAccessToken(t *testing.T) {
	accountID := int64(123)

	// 1. Probamos a generarlo
	tokenString, err := GenerateAccessToken(accountID)
	require.NoError(t, err, "No debería fallar al crear el JWT")
	require.NotEmpty(t, tokenString, "El token no puede estar vacío")

	// 2. Probamos a validarlo
	validID, err := ValidateAccessToken(tokenString)
	require.NoError(t, err, "No debería fallar al validar un JWT correcto")
	assert.Equal(t, accountID, validID, "El accountID validado debe coincidir con el original")
}

func TestValidateAccessToken_Expired(t *testing.T) {
	// Guardamos el TTL original y lo restauramos al acabar el test
	originalTTL := AccessTokenTTL
	defer func() { AccessTokenTTL = originalTTL }()

	// Forzamos un TTL negativo para generar un token ya expirado
	AccessTokenTTL = -1 * time.Second
	tokenString, _ := GenerateAccessToken(100)

	_, err := ValidateAccessToken(tokenString)
	assert.Error(t, err, "Debería devolver un error porque el token ha caducado")
}

func TestValidateAccessToken_InvalidSignature(t *testing.T) {
	// Intentar validar texto aleatorio
	_, err := ValidateAccessToken("esto.noes.unjwt")
	assert.Error(t, err, "Debería dar error al validar basura")

	// Crear un JWT pero firmado con el algoritmo incorrecto para forzar el apierror.ErrInvalidToken
	claims := jwt.MapClaims{"sub": float64(1)}
	token := jwt.NewWithClaims(jwt.SigningMethodNone, claims)
	tokenString, _ := token.SignedString(jwt.UnsafeAllowNoneSignatureType)

	_, err = ValidateAccessToken(tokenString)
	assert.Error(t, err, "Debería dar error por método de firma incorrecto")
}

func TestRefreshTokens(t *testing.T) {
	// Generar
	token1, err1 := GenerateRefreshToken()
	require.NoError(t, err1)
	require.Len(t, token1, 64, "El hex string de 32 bytes debe tener 64 caracteres de longitud")

	token2, _ := GenerateRefreshToken()
	assert.NotEqual(t, token1, token2, "Dos refresh tokens generados no pueden ser iguales")

	// Hashear
	hash1 := HashToken(token1)
	hash1Repetido := HashToken(token1)
	hash2 := HashToken(token2)

	assert.Equal(t, hash1, hash1Repetido, "El mismo token debe generar siempre el mismo hash")
	assert.NotEqual(t, hash1, hash2, "Tokens distintos deben generar hashes distintos")
	require.Len(t, hash1, 64, "El SHA256 en hex debe tener 64 caracteres")
}
