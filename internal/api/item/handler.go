package item

import (
	"net/http"
	"strconv"

	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/api/apierror"
	"github.com/gin-gonic/gin"
)

// Handler http con endpoints para tratar con items de la tienda y de usuarios

type itemHandler struct {
	service *itemService
}

func NewHandler(s *itemService) *itemHandler {
	return &itemHandler{service: s}
}

/*
*
* Desc: Esta funcion llama a otra funcion del service que obtiene la información
de un item por su id.
* --- Parametros ---
* c: *gin.Context - Es el contexto de gin, del cual recibe el id del item.
* ------------------
* Nota: si bien no hace un return de un valor, devuelve en el contexto un JSON
con la información del objeto si la obtiene junto con un StatusOK, y un error en
caso contrario.
*
*/
func (h *itemHandler) GetShopItem(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("itemID"), 10, 32)
	if err != nil {
		apierror.SendError(c, http.StatusBadRequest, err)
		return
	}

	res, err := h.service.GetByID(c.Request.Context(), int32(id))
	if err != nil {
		apierror.DetectAndSendError(c, err)
		return
	}

	c.JSON(http.StatusOK, res)
}

/*
*
* Desc: Esta funcion llama a otra funcion del service que obtiene un listado de
los objetos de un usuario dado su id.
* --- Parametros ---
* c: *gin.Context - Es el contexto de gin, del cual recibe el id del usuario.
* ------------------
* Nota: si bien no hace un return de un valor, devuelve en el contexto un JSON
con la información de los objetos si la obtiene junto con un StatusOK, y un
error en caso contrario.
*
*/
func (h *itemHandler) GetUserItems(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("userID"), 10, 64)
	if err != nil {
		apierror.SendError(c, http.StatusBadRequest, err)
		return
	}

	res, err := h.service.GetUserItems(c.Request.Context(), int64(id))
	if err != nil {
		apierror.DetectAndSendError(c, err)
		return
	}

	c.JSON(http.StatusOK, res)
}

/*
*
* Desc: Esta funcion llama a otra funcion del service que asigna un item a una
cuenta de usuario dado un JSON.
* --- Parametros ---
* c: *gin.Context - Es el contexto de gin, mediante el cual recibe el JSON.
* ------------------
* Nota: si bien no hace un return de un valor, devuelve en el contexto un JSON
con el id del usuario y del item junto con un StatusCreated si no hay errores,
y un error en caso contrario.
*
*/
func (h *itemHandler) CreateItemOwner(c *gin.Context) {

	var body ItemOwnerDTO

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
