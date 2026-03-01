package friendship

import (
	"net/http"
	"strconv"

	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/api/apierror"
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
func (h *FriendshipHandler) GetFriendship(c *gin.Context) {
	user1ID, err := strconv.ParseInt(c.Param("user1ID"), 10, 64)
	if err != nil {
		apierror.SendError(c, http.StatusBadRequest, err)
		return
	}

	user2ID, err := strconv.ParseInt(c.Param("user2ID"), 10, 64)
	if err != nil {
		apierror.SendError(c, http.StatusBadRequest, err)
		return
	}

	res, err := h.service.GetFriendshipStatus(c.Request.Context(), FriendshipDTO{
		SenderID:   user1ID,
		RecieverID: user2ID,
	})

	if err != nil {
		apierror.DetectAndSendError(c, err)
		return
	}

	c.JSON(http.StatusOK, res)
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
func (h *FriendshipHandler) GetUserFrienships(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Param("user1ID"), 10, 64)
	if err != nil {
		apierror.SendError(c, http.StatusBadRequest, err)
		return
	}

	println("El usuario es ", userID)

	res, err := h.service.GetUserFriendships(c.Request.Context(), userID)

	if err != nil {
		println("El usuario es ", userID)
		apierror.DetectAndSendError(c, err)
		return
	}

	c.JSON(http.StatusOK, res)
}

/*
*
* Desc: Esta funcion llama a otra funcion del service que crea una amistad dado
un JSON.
* --- Parametros ---
* c, *gin.Context - Es el contexto de gin de donde saca el JSON.
* ------------------
* Nota: si bien no hace un return de un valor, devuelve en el contexto un JSON
con un objeto de confirmacion que contiene los id de ambos usuarios junto con
un StatusCreated si no hay errores, y un error en caso contrario.
*
*/
func (h *FriendshipHandler) GetUserPendingSentFriendships(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Param("user1ID"), 10, 64)
	if err != nil {
		apierror.SendError(c, http.StatusBadRequest, err)
		return
	}

	res, err := h.service.GetUserPendingSentFriendships(c.Request.Context(),
		userID)

	if err != nil {
		apierror.DetectAndSendError(c, err)
		return
	}

	c.JSON(http.StatusOK, res)
}

/*
*
* Desc: Esta funcion llama a otra funcion del service que crea una amistad dado
un JSON.
* --- Parametros ---
* c, *gin.Context - Es el contexto de gin de donde saca el JSON.
* ------------------
* Nota: si bien no hace un return de un valor, devuelve en el contexto un JSON
con un objeto de confirmacion que contiene los id de ambos usuarios junto con
un StatusCreated si no hay errores, y un error en caso contrario.
*
*/
func (h *FriendshipHandler) GetUserPendingRecievedFriendships(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Param("user1ID"), 10, 64)
	if err != nil {
		apierror.SendError(c, http.StatusBadRequest, err)
		return
	}

	res, err := h.service.GetUserPendingRecievedFriendships(c.Request.Context(),
		userID)

	if err != nil {
		apierror.DetectAndSendError(c, err)
		return
	}

	c.JSON(http.StatusOK, res)
}

/*
*
* Desc: Esta funcion llama a otra funcion del service que crea una amistad dado
un JSON.
* --- Parametros ---
* c, *gin.Context - Es el contexto de gin de donde saca el JSON.
* ------------------
* Nota: si bien no hace un return de un valor, devuelve en el contexto un JSON
con un objeto de confirmacion que contiene los id de ambos usuarios junto con
un StatusCreated si no hay errores, y un error en caso contrario.
*
*/
func (h *FriendshipHandler) AcceptFrienship(c *gin.Context) {
	user1ID, err := strconv.ParseInt(c.Param("user1ID"), 10, 64)
	if err != nil {
		apierror.SendError(c, http.StatusBadRequest, err)
		return
	}

	user2ID, err := strconv.ParseInt(c.Param("user2ID"), 10, 64)
	if err != nil {
		apierror.SendError(c, http.StatusBadRequest, err)
		return
	}

	err = h.service.AcceptFrienship(c.Request.Context(), FriendshipDTO{
		SenderID:   user1ID,
		RecieverID: user2ID,
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
* Nota: si bien no hace un return de un valor, devuelve en el contexto un JSON
con un objeto de confirmacion que contiene los id de ambos usuarios junto con
un StatusCreated si no hay errores, y un error en caso contrario.
*
*/
func (h *FriendshipHandler) DeleteFrienship(c *gin.Context) {
	user1ID, err := strconv.ParseInt(c.Param("user1ID"), 10, 64)
	if err != nil {
		apierror.SendError(c, http.StatusBadRequest, err)
		return
	}

	user2ID, err := strconv.ParseInt(c.Param("user2ID"), 10, 64)
	if err != nil {
		apierror.SendError(c, http.StatusBadRequest, err)
		return
	}

	err = h.service.DeleteFrienship(c.Request.Context(), FriendshipDTO{
		SenderID:   user1ID,
		RecieverID: user2ID,
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
* Nota: si bien no hace un return de un valor, devuelve en el contexto un JSON
con un objeto de confirmacion que contiene los id de ambos usuarios junto con
un StatusCreated si no hay errores, y un error en caso contrario.
*
*/
func (h *FriendshipHandler) Create(c *gin.Context) {

	var body FriendshipDTO

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
