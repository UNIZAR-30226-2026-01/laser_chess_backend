package bot

import (
	"math/rand"
	"net/http"

	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/api/apierror"
	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/api/match"
	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/api/middleware"
	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/game"
	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/rt"
	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/sse"
	"github.com/gin-gonic/gin"
)

type BotHandler struct {
	registry     *rt.MatchRegistry
	matchService *match.MatchService
	eventSystem  *sse.EventSystem
}

func NewBotHandler(registry *rt.MatchRegistry,
	matches *match.MatchService, events *sse.EventSystem) *BotHandler {
	return &BotHandler{
		registry:     registry,
		matchService: matches,
		eventSystem:  events,
	}
}

func (h *BotHandler) BotMatch(c *gin.Context) {

	var dto BotMatchRequestDTO
	if err := c.ShouldBindQuery(&dto); err != nil {
		apierror.SendError(c, http.StatusBadRequest, err)
		return
	}

	playerID, err := middleware.ExtractAccountID(c)
	if err != nil {
		apierror.DetectAndSendError(c, err)
	}

	if _, hasMatch := h.registry.GetMatch(playerID); hasMatch {
		apierror.DetectAndSendError(c, apierror.ErrAlreadyInMatch)
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

	// Construir bot
	botClient := &rt.Client{}
	botClient.InitClient(1, nil, true)

	ai := &rt.AIClient{}
	ai.InitAIClient(botClient, game.Board_T(*dto.Board), *dto.Level)

	room := &rt.Room{}
	RedPlayer := rand.Intn(2)
	var P1Client *rt.Client
	var P2Client *rt.Client
	if RedPlayer == 0 {
		P1Client = client
		P2Client = botClient
	} else {
		P1Client = botClient
		P2Client = client
	}

	room.InitRoom(P1Client, P2Client, h.matchService, true,
		&game.GameInfo{
			BoardType:     game.Board_T(*dto.Board),
			Log:           "",
			TimeBase:      *dto.StartingTime * 1000,
			TimeIncrement: *dto.TimeIncrement * 1000,
			MatchType:     "BOTS",
		}, h.registry)

}
