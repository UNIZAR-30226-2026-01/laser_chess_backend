package public

// private_handler.go — endpoints para partidas privadas
//
// GET  api/rt/challenge        -> upgradea a WS, crea el reto y espera
// GET  api/rt/challenge/accept -> upgradea a WS, acepta el reto y arranca la Room
// GET  api/rt/challenges       -> devuelve la lista de retos pendientes (HTTP normal)

import (
	"math/rand"
	"net/http"

	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/api/account"
	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/api/apierror"
	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/api/match"
	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/api/middleware"
	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/api/rating"
	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/db/boards"
	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/game"
	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/sse"

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
	eventSystem    *sse.EventSystem
}

func NewPublicHandler(hub *rt.PublicHub, registry *rt.MatchRegistry, accounts *account.AccountService,
	matches *match.MatchService, ratings *rating.RatingService, events *sse.EventSystem) *PublicHandler {
	return &PublicHandler{
		hub:            hub,
		registry:       registry,
		accountService: accounts,
		matchService:   matches,
		ratingService:  ratings,
		gameModes: []GameMode{
			{10, 2},
			{300, 5},
			{900, 5},
			{900, 10},
			{1800, 10},
			{1800, 15},
			{3600, 15},
			{3600, 20},
		},
		eventSystem: events,
	}
}

func (h *PublicHandler) GetELOByGameMode(c *gin.Context, gameMode int,
	id int64) (*rating.RatingDTO, error) {
	var err error = nil
	var playerELO *rating.RatingDTO = nil
	switch gameMode {
	case 300:
		playerELO, err = h.ratingService.GetBlitzEloByID(c, id)
	case 900:
		playerELO, err = h.ratingService.GetRapidEloByID(c, id)
	case 1800:
		playerELO, err = h.ratingService.GetClassicEloByID(c, id)
	case 3600:
		playerELO, err = h.ratingService.GetExtendedEloByID(c, id)
	}
	if err != nil {
		return playerELO, err
	}
	return playerELO, nil
}

func (h *PublicHandler) GetIDAndELOStats(c *gin.Context, gameMode int) (int64, int32, int32) {
	// Cojer jwt
	playerID, err := middleware.ExtractAccountID(c)
	if err != nil {
		apierror.DetectAndSendError(c, err)
		return 0, 0, 0
	}

	playerELO, err := h.GetELOByGameMode(c, gameMode, playerID)
	if err != nil {
		apierror.DetectAndSendError(c, err)
		return 0, 0, 0
	}

	return playerID, playerELO.Value, playerELO.Deviation
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

	if dto.Ranked != 0 && dto.Ranked != 1 {
		apierror.SendError(c, http.StatusBadRequest, apierror.ErrNotAValidGameMode)
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

	playerID, playerELO, playerRD := h.GetIDAndELOStats(c, int(gameMode.StartingTime))

	// Comprobar que el challenger no tiene ya una partida activa
	if _, ok := h.registry.GetMatch(playerID); ok {
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
	client.InitClient(playerID, conn, false)

	// Iniciar el matchmaking

	ResponseChan := make(chan *rt.Client, 1)
	ErrorChan := make(chan error, 1)

	go h.hub.AddPlayerToMatchmaking(&rt.MatchRequest{ // TODO: anadir canal de cancel
		PlayerClient: client,
		PlayerELO:    playerELO,
		PlayerRD:     playerRD,
		GameMode:     i,
		ResponseChan: ResponseChan,
		ErrorChan:    ErrorChan,
		Ranked:       dto.Ranked,
	})

	var opponentClient *rt.Client

	select {
	case opponentClient = <-ResponseChan:
	case err := <-ErrorChan:
		apierror.DetectAndSendError(c, err)
		conn.Close()
		return
	case <-client.Done:
		// TODO: sacar de la cola
		return
	}

	if playerID > opponentClient.AccountID {
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

		var matchType string = "RANKED"
		switch dto.Ranked {
		case 0:
			matchType = "RANKED"
		case 1:
			matchType = "FRIENDLY"
		}

		room.InitRoom(P1Client, P2Client, h.matchService, true,
			&game.GameInfo{
				BoardType:     game.Board_T(rand.Intn(boards.BOARD_NUM)),
				Log:           "",
				TimeBase:      dto.StartingTime * 1000,
				TimeIncrement: dto.TimeIncrement * 1000,
				MatchType:     matchType,
			}, h.registry)

	}

}
