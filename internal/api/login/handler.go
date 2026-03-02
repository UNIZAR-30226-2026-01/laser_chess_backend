package login

import (
	"net/http"

	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/api/apierror"
	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/auth"
	"github.com/gin-gonic/gin"
)

// Handler http con endpoints para tratar con cuentas de usuario

type LoginHandler struct {
	service *LoginService
}

func NewHandler(s *LoginService) *LoginHandler {
	return &LoginHandler{service: s}
}

// Envia al cliente una cookie con el refresh token
func sendRefreshTokenCookie(c *gin.Context, res *LoginResult) {
	c.SetCookie(
		"refresh_token",                     // name
		res.RefreshToken,                    // value
		int(auth.RefreshTokenTTL.Seconds()), // maxAge
		"/",                                 // path
		"",                                  // domain
		false,                               // secure (si usamos https sera true
		true,                                // HttpOnly
	)
}

// Endpoint de login
// Devuelve un json con el access token,
// y envia una cookie con el refresh token
func (h *LoginHandler) Login(c *gin.Context) {
	var body LoginDTO
	if err := c.ShouldBindJSON(&body); err != nil {
		apierror.SendError(c, http.StatusBadRequest, err)
		return
	}

	res, err := h.service.Login(c.Request.Context(), &body)
	if err != nil {
		apierror.DetectAndSendError(c, err)
		return
	}

	// Mandar la cookie del refresh token
	sendRefreshTokenCookie(c, res)

	c.JSON(http.StatusOK, LoginResponseDTO{AccessToken: res.AccessToken})
}

// Endpoint para generar un nuevo access token
// Se presenta el refresh token y si es valido se
// devuelve un acces token nuevo
// Tambien actualiza el refresh token ya de paso
func (h *LoginHandler) Refresh(c *gin.Context) {
	// Coger el refresh de la cookie
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		apierror.SendError(c, http.StatusUnauthorized, err)
		return
	}

	// validar y rotar
	res, err := h.service.Refresh(c.Request.Context(), refreshToken)
	if err != nil {
		apierror.DetectAndSendError(c, err)
		return
	}

	// Guardar el nuevo
	sendRefreshTokenCookie(c, res)

	c.JSON(http.StatusOK, LoginResponseDTO{AccessToken: res.AccessToken})
}
