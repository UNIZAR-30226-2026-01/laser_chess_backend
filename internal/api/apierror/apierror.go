package apierror

// Paquete gestiona forma centralizada los errores para tener una interfaz común
// que utilizan los handlers para responder el código http de error correspondiente.

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

var ErrNotFound = errors.New("resource not found")

// Función que detecta el tipo de error, y manda el código de error
// correspondiente
func DetectAndSendError(c *gin.Context, err error) {
	var httpCode int
	switch {
	case errors.Is(err, ErrNotFound):
		httpCode = http.StatusNotFound
	default:
		httpCode = http.StatusInternalServerError
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
