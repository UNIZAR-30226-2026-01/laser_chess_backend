package rating

import (
	"net/http"
	"strconv"

	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/api/apierror"
	"github.com/gin-gonic/gin"
)

type RatingHandler struct {
	service *RatingService
}

func NewHandler(s *RatingService) *RatingHandler {
	return &RatingHandler{service: s}
}

func (h *RatingHandler) GetAllElos(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("userID"), 10, 64)
	if err != nil {
		apierror.DetectAndSendError(c, err)
		return
	}

	res, err := h.service.GetAllElosByID(c.Request.Context(), id)
	if err != nil {
		apierror.DetectAndSendError(c, err)
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *RatingHandler) GetBlitzElo(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("userID"), 10, 64)
	if err != nil {
		apierror.DetectAndSendError(c, err)
		return
	}

	res, err := h.service.GetBlitzEloByID(c.Request.Context(), id)
	if err != nil {
		apierror.DetectAndSendError(c, err)
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *RatingHandler) GetBulletElo(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("userID"), 10, 64)
	if err != nil {
		apierror.DetectAndSendError(c, err)
		return
	}

	res, err := h.service.GetBulletEloByID(c.Request.Context(), id)
	if err != nil {
		apierror.DetectAndSendError(c, err)
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *RatingHandler) GetRapidElo(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("userID"), 10, 64)
	if err != nil {
		apierror.DetectAndSendError(c, err)
		return
	}

	res, err := h.service.GetRapidEloByID(c.Request.Context(), id)
	if err != nil {
		apierror.DetectAndSendError(c, err)
		return
	}

	c.JSON(http.StatusOK, res)
}
func (h *RatingHandler) GetClassicElo(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("userID"), 10, 64)
	if err != nil {
		apierror.DetectAndSendError(c, err)
		return
	}

	res, err := h.service.GetClassicEloByID(c.Request.Context(), id)
	if err != nil {
		apierror.DetectAndSendError(c, err)
		return
	}

	c.JSON(http.StatusOK, res)
}

/*
func (h *RatingHandler) UpdateElo(c *gin.Context){
	var rating RatingDTO
	if err := c.ShouldBindJSON(&rating); err != nil{
		apierror.SendError(c, http.StatusBadRequest, err)
		return
	}

	res, err := h.service.UpdateEloByID(c.Request.Context(), rating)
	if err != nil {
		apierror.DetectAndSendError(c, err)
		return
	}

	c.JSON(http.StatusOK, res)
}
*/
