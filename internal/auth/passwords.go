package auth

// Fichero con funcionalidad para crear y validar contraseñas de forma segura

import "golang.org/x/crypto/bcrypt"

// Crea un hash a partir de una contraseña de texto plano
func HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// Comprueba que la contraseña provista es igual que la almacenada
func VerifyPassword(hashedPassword, providedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(providedPassword))
}
