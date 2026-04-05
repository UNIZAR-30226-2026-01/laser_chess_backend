package public

// private_handler.go — endpoints para partidas privadas
//
// GET  api/rt/challenge        -> upgradea a WS, crea el reto y espera
// GET  api/rt/challenge/accept -> upgradea a WS, acepta el reto y arranca la Room
// GET  api/rt/challenges       -> devuelve la lista de retos pendientes (HTTP normal)

import (
	"fmt"
	"math/rand"
	"net/http"

	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/api/account"
	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/api/apierror"
	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/api/match"
	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/api/middleware"
	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/api/rating"
	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/db"
	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/game"

	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/rt"
	"github.com/gin-gonic/gin"
)

type MatchmakingInfo struct {
	PlayerID      int64
	PlayerELO     int64
	GameMode      int
	StartingTime  int32
	TimeIncrement int32
}

type GameMode struct {
	StartingTime  int32
	TimeIncrement int32
}

type PublicHandler struct {
	hub            *rt.PublicHub
	registry       *rt.MatchRegistry
	accountService *account.AccountService
	matchService   *match.MatchService
	ratingService  *rating.RatingService
	gameModes      []GameMode
}

func NewPublicHandler(hub *rt.PublicHub, registry *rt.MatchRegistry, accounts *account.AccountService,
	matches *match.MatchService, ratings *rating.RatingService) *PublicHandler {
	return &PublicHandler{
		hub:            hub,
		registry:       registry,
		accountService: accounts,
		matchService:   matches,
		ratingService:  ratings,
		gameModes: []GameMode{
			{300, 2},
			{300, 5},
			{900, 5},
			{900, 10},
			{1800, 10},
			{1800, 15},
			{3600, 15},
			{3600, 20},
		},
	}
}

func (h *PublicHandler) GetIDAndElo(c *gin.Context) (int64, int64) {
	// Cojer jwt
	playerID, err := middleware.ExtractAccountID(c)
	if err != nil {
		apierror.DetectAndSendError(c, err)
		return 0, 0
	}

	playerELO, err := h.ratingService.GetBlitzEloByID(c, playerID)
	if err != nil {
		apierror.DetectAndSendError(c, err)
		return 0, 0
	}

	return playerID, int64(playerELO.Value)
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
func (h *PublicHandler) GoIntoMatchmaking(c *gin.Context) {

	var dto MatchmakingRequestDTO
	if err := c.ShouldBindQuery(&dto); err != nil {
		apierror.SendError(c, http.StatusBadRequest, err)
		return
	}

	var validGameMode bool = false
	var i int
	var gameMode GameMode
	for i, gameMode = range h.gameModes {
		validGameMode = dto.StartingTime == gameMode.StartingTime &&
			dto.TimeIncrement == gameMode.TimeIncrement

		if validGameMode {
			break
		}
	}

	if !validGameMode {
		apierror.SendError(c, http.StatusBadRequest, apierror.ErrNotAValidGameMode)
		return
	}

	playerID, playerELO := h.GetIDAndElo(c)

	request := &MatchmakingInfo{
		PlayerID:      playerID,
		PlayerELO:     playerELO,
		GameMode:      i,
		StartingTime:  dto.StartingTime,
		TimeIncrement: dto.TimeIncrement,
	}

	// Comprobar que el challenger no tiene ya una partida activa
	if _, ok := h.registry.GetMatch(request.PlayerID); ok {
		apierror.SendError(c, http.StatusConflict, apierror.ErrAlreadyInMatch)
		return
	}

	// Hacer el upgrade a websocket
	conn, err := rt.UpgradeConn(c.Writer, c.Request)
	if err != nil {
		// UpgradeConn ya manda al frontend la respuesta de error
		return
	}

	// Construir Client
	client := &rt.Client{}
	client.InitClient(request.PlayerID, conn)

	// Iniciar el matchmaking

	ResponseChan := make(chan *rt.Client, 1)

	go h.hub.AddPlayerToMatchmaking(&rt.MatchRequest{ // TODO: anadir canal de cancel
		PlayerClient: client,
		PlayerELO:    request.PlayerELO,
		GameMode:     request.GameMode,
		ResponseChan: ResponseChan,
	})
	fmt.Println("Antes del canal. Soy ", request.PlayerID)
	opponentClient := <-ResponseChan
	fmt.Println("Despues del canal. Soy ", request.PlayerID, " contra ", opponentClient.AccountID)

	if request.PlayerID > opponentClient.AccountID {
		// Creamos la partida
		room := &rt.Room{}
		RedPlayer := rand.Intn(2)
		var P1Client *rt.Client
		var P2Client *rt.Client
		if RedPlayer == 0 {
			P1Client = client
			P2Client = opponentClient
		} else {
			P1Client = opponentClient
			P2Client = client
		}
		room.InitRoom(P1Client, P2Client, h.matchService, true,
			&game.GameInfo{
				BoardType:     game.Board_T(rand.Intn(db.BOARD_NUM)),
				Log:           "",
				TimeBase:      request.StartingTime,
				TimeIncrement: request.TimeIncrement,
				MatchType:     "RANKED",
			})

		// Registramos la partida
		h.registry.RegisterMatch(request.PlayerID, opponentClient.AccountID, room)

		// Esperamos
		<-client.Done

		h.registry.RemoveMatch(request.PlayerID, opponentClient.AccountID)

	} else {
		// Esperamos
		<-opponentClient.Done

		h.registry.RemoveMatch(request.PlayerID, opponentClient.AccountID)
	}

}
