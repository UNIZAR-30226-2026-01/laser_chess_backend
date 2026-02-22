package account

import (
	"net/http"
	"strconv"

	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/api/apierror"
	db "github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/db/sqlc"
	"github.com/gin-gonic/gin"
)

// Handler http con endpoints para tratar con cuentas de usuario

type AccountHandler struct {
	service *AccountService
}

func NewHandler(s *AccountService) *AccountHandler {
	return &AccountHandler{service: s}
}

func (h *AccountHandler) CreateAccount(c *gin.Context) {
	var body db.CreateAccountParams
	// Mira si el json que nos han pasado coincide con el dto
	if err := c.ShouldBindJSON(&body); err != nil {
		apierror.SendError(c, http.StatusBadRequest, err)
		return
	}

	res, err := h.service.CreateAccount(c.Request.Context(), body)
	if err != nil {
		apierror.DetectAndSendError(c, err)
		return
	}

	c.JSON(http.StatusCreated, res)
}

func (h *AccountHandler) GetAccountByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		apierror.SendError(c, http.StatusBadRequest, err)
		return
	}

	res, err := h.service.GetAccountByID(c.Request.Context(), int64(id))
	if err != nil {
		apierror.DetectAndSendError(c, err)
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *AccountHandler) UpdateAccount(c *gin.Context) {
	var body db.UpdateAccountParams
	// Mira si el json que nos han pasado coincide con el dto
	if err := c.ShouldBindJSON(&body); err != nil {
		apierror.SendError(c, http.StatusBadRequest, err)
		return
	}

	res, err := h.service.UpdateAccount(c.Request.Context(), body)
	if err != nil {
		apierror.DetectAndSendError(c, err)
		return
	}

	c.JSON(http.StatusOK, res)
}
