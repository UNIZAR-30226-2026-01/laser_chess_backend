package reconnection

import (
	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/api/apierror"
	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/api/middleware"
	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/api/rating"
	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/rt"
	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/sse"
	"github.com/gin-gonic/gin"
)

type ReconnectionHandler struct {
	registry      *rt.MatchRegistry
	ratingService *rating.RatingService
	eventSystem   *sse.EventSystem
}

func NewReconnectionHandler(registry *rt.MatchRegistry,
	ratings *rating.RatingService, events *sse.EventSystem) *ReconnectionHandler {
	return &ReconnectionHandler{
		registry:      registry,
		ratingService: ratings,
		eventSystem:   events,
	}
}

func (h *ReconnectionHandler) Reconnect(c *gin.Context) {

	playerID, err := middleware.ExtractAccountID(c)
	if err != nil {
		apierror.DetectAndSendError(c, err)
	}

	if _, hasMatch := h.registry.GetMatch(playerID); !hasMatch {
		apierror.DetectAndSendError(c, apierror.ErrNoMatchInCourse)
	}

	// Hacer el upgrade a websocket
	conn, err := rt.UpgradeConn(c.Writer, c.Request)
	if err != nil {
		// UpgradeConn ya manda al frontend la respuesta de error
		return
	}

	// Construir Client
	client := &rt.Client{}
	client.InitClient(playerID, conn)

	if !h.registry.ReconnectClient(client) {
		apierror.DetectAndSendError(c, apierror.ErrNoMatchInCourse)
	}

}
