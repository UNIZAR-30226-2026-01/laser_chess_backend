package account

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/api/apierror"
	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/api/middleware"
	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/api/rating"
	"github.com/gin-gonic/gin"
)

// Handler http con endpoints para tratar con cuentas de usuario

type AccountHandler struct {
	accountService *AccountService
	ratingService  *rating.RatingService
}

func NewHandler(s *AccountService, r *rating.RatingService) *AccountHandler {
	return &AccountHandler{accountService: s, ratingService: r}
}

// Crea un nuevo usuario a partir de un CreateAccountDTO
func (h *AccountHandler) Create(c *gin.Context) {

	// Mira si el json que nos han pasado coincide con el dto
	var body CreateAccountDTO
	if err := c.ShouldBindJSON(&body); err != nil {
		apierror.SendError(c, http.StatusBadRequest, err)
		return
	}

	res, err := h.accountService.Create(c.Request.Context(), &body)
	if err != nil {
		apierror.DetectAndSendError(c, err)
		return
	}

	_, err = h.ratingService.CreateRating(c, *res.AccountID)
	if err != nil {
		apierror.DetectAndSendError(c, err)
		fmt.Println(err)
		return
	}

	c.JSON(http.StatusCreated, res)
}

// Devuelve un AccountDTO lleno con toda la info
// de tu propia cuenta.
func (h *AccountHandler) GetOwnAccount(c *gin.Context) {
	accountID, err := middleware.ExtractAccountID(c)
	if err != nil {
		apierror.DetectAndSendError(c, err)
		return
	}

	res, err := h.accountService.GetByID(c.Request.Context(), int64(accountID))
	if err != nil {
		apierror.DetectAndSendError(c, err)
		return
	}

	c.JSON(http.StatusOK, res)
}

// Devuelve un AccountDTO con la info publica
// de una cuenta.
func (h *AccountHandler) GetOtherByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		apierror.SendError(c, http.StatusBadRequest, err)
		return
	}

	res, err := h.accountService.GetByID(c.Request.Context(), int64(id))
	if err != nil {
		apierror.DetectAndSendError(c, err)
		return
	}

	res.Mail = nil
	res.Money = nil

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

	res, err := h.accountService.Update(c.Request.Context(), accountID, &body)
	if err != nil {
		apierror.DetectAndSendError(c, err)
		return
	}

	c.JSON(http.StatusOK, res)
}

// Desactiva la cuenta del user que manda la peticion
func (h *AccountHandler) Delete(c *gin.Context) {

	id, err := middleware.ExtractAccountID(c)

	err = h.accountService.Delete(c.Request.Context(), int64(id))
	if err != nil {
		apierror.DetectAndSendError(c, err)
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}
