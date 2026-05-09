package device

import (
	"fmt"
	"net/http"

	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/api/apierror"
	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/api/middleware"
	"github.com/gin-gonic/gin"
)

// Handler http con endpoints para tratar con cuentas de usuario

type DeviceHandler struct {
	deviceService *DeviceService
}

func NewHandler(s *DeviceService) *DeviceHandler {
	return &DeviceHandler{deviceService: s}
}

// Crea un nuevo usuario a partir de un CreateAccountDTO
func (h *DeviceHandler) RegisterDevice(c *gin.Context) {

	accountID, err := middleware.ExtractAccountID(c)
	if err != nil {
		apierror.DetectAndSendError(c, err)
		return
	}

	// Mira si el json que nos han pasado coincide con el dto
	var body DeviceDTO
	if err := c.ShouldBindJSON(&body); err != nil {
		apierror.SendError(c, http.StatusBadRequest, err)
		return
	}

	fmt.Println(body)

	res, err := h.deviceService.RegisterDevice(c.Request.Context(),
		body, accountID)

	fmt.Println(err.Error())
	if err != nil {
		apierror.DetectAndSendError(c, err)
		return
	}

	c.JSON(http.StatusCreated, res)
}

func (h *DeviceHandler) DeleteDevice(c *gin.Context) {

	_, err := middleware.ExtractAccountID(c)
	if err != nil {
		apierror.DetectAndSendError(c, err)
		return
	}

	// Mira si el json que nos han pasado coincide con el dto
	var body DeviceDTO
	if err := c.ShouldBindJSON(&body); err != nil {
		apierror.SendError(c, http.StatusBadRequest, err)
		return
	}

	fmt.Println(body.Token)
	res, err := h.deviceService.DeleteDevice(c.Request.Context(),
		body.Token)

	if err != nil {
		apierror.DetectAndSendError(c, err)
		return
	}

	c.JSON(http.StatusCreated, res)
}
