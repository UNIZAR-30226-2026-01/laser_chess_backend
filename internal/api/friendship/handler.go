package friendship

import (
	"net/http"
	"strconv"

	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/api/apierror"
	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/api/middleware"
	"github.com/gin-gonic/gin"
)

type FriendshipHandler struct {
	service *FriendshipService
}

func NewHandler(s *FriendshipService) *FriendshipHandler {
	return &FriendshipHandler{service: s}
}

/*
*
* Desc: Esta funcion llama a otra funcion del service que busca una amistad
dados los ids de los usuarios.
* --- Parametros ---
* c, *gin.Context - Es el contexto de gin del cual obtiene los ids.
* ------------------
* Nota: si bien no hace un return de un valor, devuelve en el contexto un JSON
con la informacion de la amistad junto con un StatusOK si no hay errores, y un
error en caso contrario.
*
*/
func (h *FriendshipHandler) GetFriendshipStatus(c *gin.Context) {
	user1ID, err := middleware.ExtractAccountID(c)
	if err != nil {
		apierror.DetectAndSendError(c, err)
		return
	}

	user2ID, err := strconv.ParseInt(c.Param("user2ID"), 10, 64)
	if err != nil {
		apierror.SendError(c, http.StatusBadRequest, err)
		return
	}

	res, err := h.service.GetFriendshipStatus(c.Request.Context(), &FriendshipDTO{
		SenderID:   &user1ID,
		ReceiverID: user2ID,
	})

	if err != nil {
		apierror.DetectAndSendError(c, err)
		return
	}

	c.JSON(http.StatusOK, res)
}

/*
*
* Desc: Esta funcion llama a otra funcion del service que devuelve un listado de
las amistades de un usuario dado su id.
* --- Parametros ---
* c, *gin.Context - Es el contexto de gin del cual obtiene el id del usuario.
* ------------------
* Nota: si bien no hace un return de un valor, devuelve en el contexto un JSON
con el listado de las amistades que contienen los datos relevantes del otro
usuario junto con un StatusOK si no hay errores, y un error en caso contrario.
*
*/
func (h *FriendshipHandler) GetUserFriendships(c *gin.Context) {
	accountID, err := middleware.ExtractAccountID(c)
	if err != nil {
		apierror.DetectAndSendError(c, err)
		return
	}

	res, err := h.service.GetUserFriendships(c.Request.Context(), accountID)
	if err != nil {
		apierror.DetectAndSendError(c, err)
		return
	}

	c.JSON(http.StatusOK, res)
}

/*
*
* Desc: Esta funcion llama a otra funcion del service que devuelve un listado de
las amistades enviadas pendientes de un usuario dado su id.
* --- Parametros ---
* c, *gin.Context - Es el contexto de gin del cual obtiene el id del usuario.
* ------------------
* Nota: si bien no hace un return de un valor, devuelve en el contexto un JSON
con el listado de las amistades enviadas pendientes que contienen los datos
relevantes del otro usuario junto con un StatusOK si no hay errores, y un error
en caso contrario.
*
*/
func (h *FriendshipHandler) GetUserPendingSentFriendships(c *gin.Context) {
	accountID, err := middleware.ExtractAccountID(c)
	if err != nil {
		apierror.DetectAndSendError(c, err)
		return
	}

	res, err := h.service.GetUserPendingSentFriendships(c.Request.Context(), accountID)

	if err != nil {
		apierror.DetectAndSendError(c, err)
		return
	}

	c.JSON(http.StatusOK, res)
}

/*
*
* Desc: Esta funcion llama a otra funcion del service que devuelve un listado de
las amistades recibidas pendientes de un usuario dado su id.
* --- Parametros ---
* c, *gin.Context - Es el contexto de gin del cual obtiene el id del usuario.
* ------------------
* Nota: si bien no hace un return de un valor, devuelve en el contexto un JSON
con el listado de las amistades recibidas pendientes que contienen los datos
relevantes del otro usuario junto con un StatusOK si no hay errores, y un error
en caso contrario.
*
*/
func (h *FriendshipHandler) GetUserPendingReceivedFriendships(c *gin.Context) {
	accountID, err := middleware.ExtractAccountID(c)
	if err != nil {
		apierror.DetectAndSendError(c, err)
		return
	}

	res, err := h.service.GetUserPendingReceivedFriendships(c.Request.Context(), accountID)

	if err != nil {
		apierror.DetectAndSendError(c, err)
		return
	}

	c.JSON(http.StatusOK, res)
}

/*
*
* Desc: Esta funcion llama a otra funcion del service que acepta una amistad
entre dos usuarios dados sus ids.
* --- Parametros ---
* c, *gin.Context - Es el contexto de gin de donde los ids de los usuarios.
* ------------------
* Nota: si bien no hace un return de un valor, devuelve en el contexto un JSON
un error en caso de haber ocurrido.
*
*/
func (h *FriendshipHandler) AcceptFriendship(c *gin.Context) {
	user1ID, err := middleware.ExtractAccountID(c)
	if err != nil {
		apierror.DetectAndSendError(c, err)
		return
	}

	user2ID, err := strconv.ParseInt(c.Param("user2ID"), 10, 64)
	if err != nil {
		apierror.SendError(c, http.StatusBadRequest, err)
		return
	}

	err = h.service.AcceptFriendship(c.Request.Context(), &FriendshipDTO{
		SenderID:   &user1ID,
		ReceiverID: user2ID,
	})

	if err != nil {
		apierror.DetectAndSendError(c, err)
		return
	}

	c.JSON(http.StatusOK, nil)
}

/*
*
* Desc: Esta funcion llama a otra funcion del service que elimina una amistad
entre dos usuarios dados sus ids.
* --- Parametros ---
* c, *gin.Context - Es el contexto de gin de donde los ids de los usuarios.
* ------------------
* Nota: si bien no hace un return de un valor, devuelve en el contexto un JSON
un error en caso de haber ocurrido.
*
*/
func (h *FriendshipHandler) DeleteFriendship(c *gin.Context) {
	user1ID, err := middleware.ExtractAccountID(c)
	if err != nil {
		apierror.DetectAndSendError(c, err)
		return
	}

	user2ID, err := strconv.ParseInt(c.Param("user2ID"), 10, 64)
	if err != nil {
		apierror.SendError(c, http.StatusBadRequest, err)
		return
	}

	err = h.service.DeleteFriendship(c.Request.Context(), &FriendshipDTO{
		SenderID:   &user1ID,
		ReceiverID: user2ID,
	})

	if err != nil {
		apierror.DetectAndSendError(c, err)
		return
	}

	c.JSON(http.StatusOK, nil)
}

/*
*
* Desc: Esta funcion llama a otra funcion del service que crea una amistad dado
un JSON.
* --- Parametros ---
* c, *gin.Context - Es el contexto de gin de donde saca el JSON.
* ------------------
Devuelve StatusCreated si no hay errores, y un error en caso contrario.
*
*/
func (h *FriendshipHandler) Create(c *gin.Context) {

	accountID, err := middleware.ExtractAccountID(c)
	if err != nil {
		apierror.DetectAndSendError(c, err)
		return
	}

	var body FriendshipDTO

	// Mira si el json que nos han pasado coincide con el dto
	if err := c.ShouldBindJSON(&body); err != nil {
		apierror.SendError(c, http.StatusBadRequest, err)
		return
	}

	body.SenderID = &accountID
	err = h.service.Create(c.Request.Context(), &body)
	if err != nil {
		apierror.DetectAndSendError(c, err)
		return
	}

	c.JSON(http.StatusCreated, nil)
}
