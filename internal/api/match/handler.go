package match

import (
	"net/http"
	"strconv"

	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/api/apierror"
	db "github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/db/sqlc"
	"github.com/gin-gonic/gin"
)

type MatchHandler struct {
	service *MatchService
}

func NewHandler(s *MatchService) *MatchHandler {
	return &MatchHandler{service: s}
}

func (h *MatchHandler) GetMatch(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("matchID"), 10, 32)
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

func (h *MatchHandler) GetUserHistory(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("userID"), 10, 32)
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

func (h *MatchHandler) CreateMatch(c *gin.Context) {

	var body db.CreateMatchParams

	// Mira si el json que nos han pasado coincide con el dto
	if err := c.ShouldBindJSON(&body); err != nil {
		apierror.SendError(c, http.StatusBadRequest, err)
		return
	}

	res, err := h.service.Create(c.Request.Context(), body)
	if err != nil {
		apierror.DetectAndSendError(c, err)
		return
	}

	c.JSON(http.StatusCreated, res)
}
