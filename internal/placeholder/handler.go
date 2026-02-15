package placeholder

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PlaceholderHandler struct {
	service *PlaceholderService
}

func NewHandler(s *PlaceholderService) *PlaceholderHandler {
	return &PlaceholderHandler{service: s}
}

func (h *PlaceholderHandler) GetPlaceholder(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		h.sendError(c, http.StatusBadRequest, err)
		return
	}

	res, err := h.service.GetByID(c.Request.Context(), int32(id))
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			h.sendError(c, http.StatusNotFound, err)
		} else {
			h.sendError(c, http.StatusInternalServerError, err)
		}
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *PlaceholderHandler) CreatePlaceholder(c *gin.Context) {
	// Struct privado con el json que nos van a pasar
	// Es básicamente un dto que solo se usa aqui
	// Probablemente creemos algunos paquetes con dtos
	var body struct {
		Data string `json:"data" binding:"required"`
	}

	// Mira si el json que nos han pasado coincide con el dto
	if err := c.ShouldBindJSON(&body); err != nil {
		h.sendError(c, http.StatusBadRequest, err)
		return
	}

	res, err := h.service.Create(c.Request.Context(), body.Data)
	if err != nil {
		h.sendError(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, res)
}

// Método auxiliar para enviar errores genéricos
// Habrá que meterlo en algun paquete de errores
func (h *PlaceholderHandler) sendError(c *gin.Context, code int, err error) {
	if err != nil {
		log.Printf("DEBUG ERROR [%d]: %v", code, err)
	}
	c.AbortWithStatusJSON(code, gin.H{
		"error": http.StatusText(code),
	})
}
