package friendship

import (
	"net/http"
	"strconv"

	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/api/apierror"
	"github.com/gin-gonic/gin"
)

type FriendshipHandler struct {
	service *friendshipService
}

func NewHandler(s *friendshipService) *friendshipHandler {
	return &friendshipHandler{service: s}
}

/**
* Desc: Esta función devuelve una lista de amistades consolidadas con otros 
*		usuarios y gestiona que las llamadas se hagan correctamente.
* --- parametros ---
* «h» friendshipHandler : "objeto handler" encargado de gestionar la llamada a BDD
* «c» *gin.Context : petición c.Param y contexto de gin
* --- resultados ---
* «c» *gin.Context: resultado c.JSON, con estado de petición http y resultado
*/
func (h *friendshipHandler) getUserFriendships(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("userID", 10, 32))
	if err != nil {
		apierror.SendError(c, http.StatusBadRequest, err)
		return
	}

	/*TODO : PETICIÓN A BDD*/

	c.JSON(http.StatusOK, res)
}

