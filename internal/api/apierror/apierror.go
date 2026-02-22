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
	ErrNotFound      = errors.New("resource not found")
	ErrAlreadyExists = errors.New("resource already exists")
	ErrBadRequest    = errors.New("bad request")
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

		case pgerrcode.ForeignKeyViolation, pgerrcode.NotNullViolation,
			pgerrcode.CheckViolation:

			httpCode = http.StatusBadRequest
		default:
			// Cualquier otro error raro de BD
			httpCode = http.StatusInternalServerError
		}
	} else {
		switch {
		case errors.Is(err, pgx.ErrNoRows) || errors.Is(err, ErrNotFound):
			httpCode = http.StatusNotFound

		case errors.Is(err, ErrAlreadyExists):
			httpCode = http.StatusConflict

		case errors.Is(err, ErrBadRequest):
			httpCode = http.StatusBadRequest

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
		"error": http.StatusText(code),
	})
}
