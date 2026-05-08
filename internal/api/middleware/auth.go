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
func AuthMiddleware(wsStore *auth.WsTicketStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		var accountID int64
		var err error

		header := c.GetHeader("Authorization")

		if header == "" {
			header = c.Query("token")
			header = "Bearer " + header
		}

		if header != "" && strings.HasPrefix(header, "Bearer ") && header != "Bearer " {
			tokenString := strings.TrimPrefix(header, "Bearer ")
			accountID, err = auth.ValidateAccessToken(tokenString)
			if err != nil {
				apierror.DetectAndSendError(c, apierror.ErrInvalidToken)
				c.Abort()
				return
			}
		} else {
			ticket := c.Query("ticket")
			if ticket != "" {
				var isValid bool
				accountID, isValid = wsStore.GetAndDelete(ticket)
				if !isValid {
					apierror.DetectAndSendError(c, apierror.ErrInvalidToken)
					c.Abort()
					return
				}
			} else {
				// No hay ni JWT válido ni Ticket
				apierror.DetectAndSendError(c, apierror.ErrInvalidToken)
				c.Abort()
				return
			}
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
