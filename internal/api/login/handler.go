package login

import (
	"net/http"

	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/api/apierror"
	"github.com/gin-gonic/gin"
)

// Handler http con endpoints para tratar con cuentas de usuario

type LoginHandler struct {
	service *LoginService
}

func NewHandler(s *LoginService) *LoginHandler {
	return &LoginHandler{service: s}
}

// Endpoint de login
func (h *LoginHandler) Login(c *gin.Context) {
	var body LoginDTO
	if err := c.ShouldBindJSON(&body); err != nil {
		apierror.SendError(c, http.StatusBadRequest, err)
		return
	}

	res, err := h.service.Login(c.Request.Context(), body)
	if err != nil {
		apierror.DetectAndSendError(c, err)
		return
	}

	// Mandar la cookie del refresh token
	c.SetCookie(
		"refresh_token",  // name
		res.RefreshToken, // value
		3600*24*7,        // maxAge
		"/",              // path
		"",               // domain
		false,            // secure (si usamos https sera true
		true,             // HttpOnly
	)

	c.JSON(http.StatusOK, LoginResponseDTO{AccessToken: res.AccessToken})
}
