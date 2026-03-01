package match

import (
	"net/http"
	"strconv"

	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/api/apierror"
	"github.com/gin-gonic/gin"
)

type MatchHandler struct {
	service *MatchService
}

func NewHandler(s *MatchService) *MatchHandler {
	return &MatchHandler{service: s}
}

/*
*
* Desc: Esta funcion llama a otra funcion del service que busca una partida dado
su id.
* --- Parametros ---
* c, *gin.Context - Es el contexto de gin del cual obtiene el id de la partida.
* ------------------
* Nota: si bien no hace un return de un valor, devuelve en el contexto un JSON
con la informacion de la partida junto con un StatusOK si no hay errores, y un
error en caso contrario.
*
*/
func (h *MatchHandler) GetMatch(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("matchID"), 10, 64)
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

/*
*
* Desc: Esta funcion llama a otra funcion del service que busca el listado de
partidas de un usuario dado su id.
* --- Parametros ---
* c, *gin.Context - Es el contexto de gin del cual obtiene el id del usuario.
* ------------------
* Nota: si bien no hace un return de un valor, devuelve en el contexto un JSON
con la lista de partidas del usuario junto con un StatusOK si no hay errores, y
un error en caso contrario.
*
*/
func (h *MatchHandler) GetUserHistory(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("userID"), 10, 64)
	if err != nil {
		apierror.SendError(c, http.StatusBadRequest, err)
		return
	}

	res, err := h.service.GetUserHistory(c.Request.Context(), int64(id))
	if err != nil {
		apierror.DetectAndSendError(c, err)
		return
	}

	c.JSON(http.StatusOK, res)
}

/*
*
* Desc: Esta funcion llama a otra funcion del service que crea una partida dado
un JSON.
* --- Parametros ---
* c, *gin.Context - Es el contexto de gin de donde saca el JSON.
* ------------------
* Nota: si bien no hace un return de un valor, devuelve en el contexto un JSON
con un objeto de confirmacion que contiene los id de ambos jugadores junto con
un StatusCreated si no hay errores, y un error en caso contrario.
*
*/
func (h *MatchHandler) CreateMatch(c *gin.Context) {

	var body MatchDTO

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
