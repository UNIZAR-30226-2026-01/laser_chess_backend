package rt

// private_handler.go — endpoints para partidas privadas
//
// GET  api/rt/challenge        -> upgradea a WS, crea el reto y espera
// GET  api/rt/challenge/accept -> upgradea a WS, acepta el reto y arranca la Room
// GET  api/rt/challenges       -> devuelve la lista de retos pendientes (HTTP normal)

import (
	"net/http"

	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/api/account"
	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/api/apierror"
	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/api/middleware"
	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/game"
	"github.com/gin-gonic/gin"
)

type PrivateHandler struct {
	hub            *PrivateHub
	registry       *MatchRegistry
	accountService *account.AccountService
}

func NewPrivateHandler(hub *PrivateHub, registry *MatchRegistry, accounts *account.AccountService) *PrivateHandler {
	return &PrivateHandler{
		hub:            hub,
		registry:       registry,
		accountService: accounts,
	}
}

// Challenge — el challenger upgradea a WS y queda a la espera de que lo acepten.
//
// Query params: username, board, starting_time, time_increment
//
//	Valida params y saca el challengerID del JWT.
//	Consigue el ID del challenger a partir del username.
//	Upgradea la conexión a WebSocket.
//	Registra el reto en el PrivateHub.
//	Bloquea leyendo del canal Send del cliente (mensajes de la Room
//	   cuando la partida arranque) o hasta que el WS se cierre.
func (h *PrivateHandler) Challenge(c *gin.Context) {
	// Coger params y jwt
	var dto CreateChallengeDTO
	if err := c.ShouldBindQuery(&dto); err != nil {
		apierror.SendError(c, http.StatusBadRequest, err)
		return
	}

	challengerID, err := middleware.ExtractAccountID(c)
	if err != nil {
		apierror.DetectAndSendError(c, err)
		return
	}

	// Conseguir id del challenged
	challengedID, err := h.accountService.GetIDByUsername(c.Request.Context(), dto.ChallengedUsername)
	if err != nil {
		apierror.DetectAndSendError(c, err)
		return
	}

	// No puedes retarte a ti mismo
	if challengerID == challengedID {
		apierror.SendError(c, http.StatusBadRequest, apierror.ErrSelfChallenge)
		return
	}

	// Comprobar que el challenger no tiene ya una partida activa
	if _, ok := h.registry.GetMatch(challengerID); ok {
		apierror.SendError(c, http.StatusConflict, apierror.ErrAlreadyInMatch)
		return
	}

	// Hacer el upgrade a websocket
	conn, err := UpgradeConn(c.Writer, c.Request)
	if err != nil {
		// UpgradeConn ya manda al frontend la respuesta de error
		return
	}

	// Construir Client (sin Room todavia)
	toRoom := make(chan interface{})
	client := &Client{}
	client.InitClient(challengerID, conn, toRoom)

	// Registrar reto en el hub privado
	info := &ChallengeInfo{
		challengerClient: client,
		challengedId:     challengedID,
		board:            game.Board_T(dto.Board),
		startingTime:     dto.StartingTime,
		timeIncrement:    dto.TimeIncrement,
	}

	err = h.hub.CreateChallenge(challengerID, challengedID, info)
	if err != nil {
		client.Close()
		apierror.SendError(c, http.StatusConflict, err)
		return
	}

	// Esperar a que el WS se cierre.
	// Si el challenger cancela antes de que lo acepten, limpiamos el reto.
	// Si la partida arranca, la Room cerrará la conn al terminar.
	<-client.Done
	h.hub.RemoveChallenge(challengerID, challengedID)

}

// AcceptChallenge — el challenged upgradea a WS y acepta el reto.
//
// Query params: username  (username del challenger)
//
//	Valida params y saca el challengedID del JWT.
//	Consigue el ID del challenger a partir del username.
//	Upgradea a WebSocket.
//	Llama a AcceptChallenge en el hub → recibe ChallengeInfo.
//	Crea la Room con ambos clientes y arranca la partida.
func (h *PrivateHandler) AcceptChallenge(c *gin.Context) {

	// Cojer params y JWT
	var dto AcceptChallengeDTO
	if err := c.ShouldBindQuery(&dto); err != nil {
		apierror.SendError(c, http.StatusBadRequest, err)
		return
	}

	challengedID, err := middleware.ExtractAccountID(c)
	if err != nil {
		apierror.DetectAndSendError(c, err)
		return
	}

	// Conseguir id del challenged
	challengerID, err := h.accountService.GetIDByUsername(c.Request.Context(), dto.ChallengerUsername)
	if err != nil {
		apierror.DetectAndSendError(c, err)
		return
	}

	// Comprobar que ninguno de los dos ya está en partida
	if _, ok := h.registry.GetMatch(challengedID); ok {
		apierror.SendError(c, http.StatusConflict, apierror.ErrAlreadyInMatch)
		return
	}
	if _, ok := h.registry.GetMatch(challengerID); ok {
		apierror.SendError(c, http.StatusConflict, apierror.ErrAlreadyInMatch)
		return
	}

	// Upgrade a WebSocket
	conn, err := UpgradeConn(c.Writer, c.Request)
	if err != nil {
		return
	}

	// Aceptar el reto en el hub
	info, err := h.hub.AcceptChallenge(challengerID, challengedID)
	if err != nil {
		conn.Close()
		apierror.SendError(c, http.StatusNotFound, err)
		return
	}

	// Construir el Client del challenged
	toRoom := make(chan interface{})
	challengedClient := &Client{}
	challengedClient.InitClient(challengedID, conn, toRoom)

	// Crear la Room y arrancar la partida
	room := &Room{}
	room.InitRoom(info.challengerClient, challengedClient, game.ACE) // ESPECIFICAR EL TABLERO

	// Registrar ambos jugadores en el registry
	h.registry.RegisterMatch(challengerID, challengedID, room)

}

// GetChallenges — devuelve la lista de retos pendientes recibidos por el usuario.
//
// Response: []PendingChallengeDTO
func (h *PrivateHandler) GetChallenges(c *gin.Context) {
	accountID, err := middleware.ExtractAccountID(c)
	if err != nil {
		apierror.DetectAndSendError(c, err)
		return
	}

	challengerIDs := h.hub.GetChallenges(accountID)
	if len(challengerIDs) == 0 {
		c.JSON(http.StatusOK, []PendingChallengeDTO{})
		return
	}

	result := make([]PendingChallengeDTO, 0, len(challengerIDs))
	for _, cID := range challengerIDs {
		info := h.hub.GetChallengeInfo(cID)
		if info == nil {
			continue
		}

		username, err := h.accountService.GetUsernameByID(c.Request.Context(), cID)
		if err != nil {
			continue
		}

		result = append(result, PendingChallengeDTO{
			ChallengerID:       cID,
			ChallengerUsername: username,
			Board:              int(info.board),
			StartingTime:       info.startingTime,
			TimeIncrement:      info.timeIncrement,
		})
	}

	c.JSON(http.StatusOK, result)
}
