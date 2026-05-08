package device

import (
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
	var body RegisterDeviceDTO
	if err := c.ShouldBindJSON(&body); err != nil {
		apierror.SendError(c, http.StatusBadRequest, err)
		return
	}

	res, err := h.deviceService.RegisterDevice(c.Request.Context(),
		body, accountID)

	if err != nil {
		apierror.DetectAndSendError(c, err)
		return
	}

	c.JSON(http.StatusCreated, res)
}
