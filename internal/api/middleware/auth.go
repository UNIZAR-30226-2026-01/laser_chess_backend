package middleware

// Middleware que intercepta peticiones http a endpoints protegidos
// y solo deja pasar si hay un token valido

import (
	"strings"

	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/api/apierror"
	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/auth"
	"github.com/gin-gonic/gin"
)

// Middleware que intercepta peticiones http a endpoints protegidos
// y solo deja pasar si hay un token valido
// Setea el userID en el context para que lo puedan usar los handlers
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			header = c.Query("token")
			header = "Bearer " + header
		}

		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			error := apierror.ErrInvalidToken
			apierror.DetectAndSendError(c, error)
			return
		}

		tokenString := strings.TrimPrefix(header, "Bearer ")
		accountID, err := auth.ValidateAccessToken(tokenString)
		if err != nil {
			error := apierror.ErrInvalidToken
			apierror.DetectAndSendError(c, error)
			return
		}

		// Guardamos el ID para que los handlers sepan quién es el usuario
		c.Set("account_id", accountID)
		c.Next()
	}
}

// Extrae el id de la cuenta que envia la peticion http
// con access token jwt
func ExtractAccountID(c *gin.Context) (int64, error) {

	idInterface, exists := c.Get("account_id")
	if !exists {
		return -1, apierror.ErrInternalServerError
	}

	accountID, ok := idInterface.(int64)
	if !ok {
		return -1, apierror.ErrInternalServerError
	}

	return accountID, nil
}
