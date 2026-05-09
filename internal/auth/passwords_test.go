package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestHashAndVerifyPassword(t *testing.T) {
	password := "SecretoAjedrezLaser123!"

	// Generar hash
	hash, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hash)

	// Validar correcto
	errMatch := VerifyPassword(hash, password)
	assert.NoError(t, errMatch)

	// Validar incorrecto
	wrongPassword := "OtraCosa!"
	errMismatch := VerifyPassword(hash, wrongPassword)
	assert.ErrorIs(t, errMismatch, bcrypt.ErrMismatchedHashAndPassword)
}
