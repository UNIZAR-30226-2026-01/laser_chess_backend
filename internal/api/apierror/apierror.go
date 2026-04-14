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
	ErrInvalidToken        = errors.New("invalid token")
	ErrUnauthorized        = errors.New("unauthorized access")

	ErrInvalidPasswordLenght = errors.New("password must be between 6 and 50 characters")
	ErrInvalidMailFormat     = errors.New("mail format is invalid")

	ErrSelfChallenge        = errors.New("you can't challenge yourself")
	ErrNotFriends           = errors.New("users are not friends")
	ErrAlreadyInMatch       = errors.New("user already in match")
	ErrAlreadyInQueue       = errors.New("user already in queue")
	ErrMatchAlreadyFinished = errors.New("match is already finished")
	ErrNotYourMatch         = errors.New("match is not yours")
	ErrNotAValidGameMode    = errors.New("time base and/or increment invalid")
	ErrNotAValidRankedType  = errors.New("ranked type invalid")
	ErrNoMatchInCourse      = errors.New("user doesn't have any running matches")

	ErrAlreadyFriends = errors.New("users are already friends")

	ErrNotEnoughMoney  = errors.New("user doesn't have enough money")
	ErrUserLevelTooLow = errors.New("user's level doesn't meet the level requisite")
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
		case errors.Is(err, pgx.ErrNoRows) || errors.Is(err, ErrNotFound):
			httpCode = http.StatusNotFound
			err = ErrNotFound

		case errors.Is(err, ErrAlreadyExists):
			httpCode = http.StatusConflict

		case errors.Is(err, ErrUnauthorized):
			httpCode = http.StatusUnauthorized

		case errors.Is(err, ErrBadRequest):
			httpCode = http.StatusBadRequest

		case errors.Is(err, ErrInvalidToken):
			httpCode = http.StatusUnauthorized

		case errors.Is(err, ErrNoMatchInCourse):
			httpCode = http.StatusNotFound

		case errors.Is(err, ErrMatchAlreadyFinished):
			httpCode = http.StatusBadRequest

		case errors.Is(err, ErrNotYourMatch):
			httpCode = http.StatusBadRequest

		case errors.Is(err, ErrSelfChallenge):
			httpCode = http.StatusBadRequest

		case errors.Is(err, ErrNotFriends):
			httpCode = http.StatusBadRequest

		case errors.Is(err, ErrAlreadyInMatch):
			httpCode = http.StatusConflict

		case errors.Is(err, ErrAlreadyInQueue):
			httpCode = http.StatusConflict

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
