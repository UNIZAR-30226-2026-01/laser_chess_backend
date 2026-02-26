package apierror

// Paquete gestiona forma centralizada los errores para tener una interfaz común
// que utilizan los handlers para responder el código http de error correspondiente.

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrInternalServerError = errors.New("internal server error")
	ErrNotFound            = errors.New("resource not found")
	ErrAlreadyExists       = errors.New("resource already exists")
	ErrBadRequest          = errors.New("bad request")
	ErrUnauthorized        = errors.New("unauthorized access")
)

// Función que detecta el tipo de error, y manda el código de error
// correspondiente
func DetectAndSendError(c *gin.Context, err error) {
	var httpCode int

	// primero miramos errores de postgres
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case pgerrcode.UniqueViolation:
			httpCode = http.StatusConflict
			err = ErrInternalServerError

		case pgerrcode.ForeignKeyViolation, pgerrcode.NotNullViolation,
			pgerrcode.CheckViolation:

			httpCode = http.StatusBadRequest
			err = ErrInternalServerError
		default:
			// Cualquier otro error raro de BD
			httpCode = http.StatusInternalServerError
			err = ErrInternalServerError
		}
	} else {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			httpCode = http.StatusNotFound
			err = ErrNotFound

		case errors.Is(err, ErrUnauthorized):
			httpCode = http.StatusUnauthorized
			err = ErrUnauthorized

		default:
			httpCode = http.StatusInternalServerError
		}
	}

	SendError(c, httpCode, err)
}

// Método auxiliar para enviar errores
func SendError(c *gin.Context, code int, err error) {
	if err != nil {
		log.Printf("DEBUG ERROR [%d]: %v", code, err)
	}
	c.AbortWithStatusJSON(code, gin.H{
		"error":         http.StatusText(code),
		"error_message": err.Error(),
	})
}
