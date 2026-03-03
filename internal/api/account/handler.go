package account

import (
	"net/http"
	"strconv"

	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/api/apierror"
	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/api/middleware"
	"github.com/gin-gonic/gin"
)

// Handler http con endpoints para tratar con cuentas de usuario

type AccountHandler struct {
	service *AccountService
}

func NewHandler(s *AccountService) *AccountHandler {
	return &AccountHandler{service: s}
}

// Crea un nuevo usuario a partir de un CreateAccountDTO
func (h *AccountHandler) Create(c *gin.Context) {

	// Mira si el json que nos han pasado coincide con el dto
	var body CreateAccountDTO
	if err := c.ShouldBindJSON(&body); err != nil {
		apierror.SendError(c, http.StatusBadRequest, err)
		return
	}

	res, err := h.service.Create(c.Request.Context(), &body)
	if err != nil {
		apierror.DetectAndSendError(c, err)
		return
	}

	c.JSON(http.StatusCreated, res)
}

// Devuelve un AccountDTO lleno con toda la info
// de una cuenta. El id de cuenta se pasa en la url
func (h *AccountHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		apierror.SendError(c, http.StatusBadRequest, err)
		return
	}

	res, err := h.service.GetByID(c.Request.Context(), int64(id))
	if err != nil {
		apierror.DetectAndSendError(c, err)
		return
	}

	c.JSON(http.StatusOK, res)
}

// Actualiza el username y/o los cosmeticos equipados de un user
// Solo hace falta mandar los campos que se vayan a actualizar,
// el resto pueden estar vacios
// El account id lo pilla del access jwt
func (h *AccountHandler) Update(c *gin.Context) {
	accountID, err := middleware.ExtractAccountID(c)
	if err != nil {
		apierror.DetectAndSendError(c, err)
		return
	}

	// Mira si el json que nos han pasado coincide con el dto
	var body AccountDTO
	if err := c.ShouldBindJSON(&body); err != nil {
		apierror.SendError(c, http.StatusBadRequest, err)
		return
	}

	res, err := h.service.Update(c.Request.Context(), accountID, &body)
	if err != nil {
		apierror.DetectAndSendError(c, err)
		return
	}

	c.JSON(http.StatusOK, res)
}

// Desactiva la cuenta del user que manda la peticion
func (h *AccountHandler) Delete(c *gin.Context) {

	id, err := middleware.ExtractAccountID(c)

	err = h.service.Delete(c.Request.Context(), int64(id))
	if err != nil {
		apierror.DetectAndSendError(c, err)
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}
