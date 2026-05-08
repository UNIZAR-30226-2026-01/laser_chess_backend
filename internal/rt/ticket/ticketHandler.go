package ticket

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"time"

	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/api/apierror"
	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/api/middleware"
	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/auth"
	"github.com/gin-gonic/gin"
)

type TicketHandler struct {
	ticketStore *auth.WsTicketStore
}

func NewHandler(ts *auth.WsTicketStore) *TicketHandler {
	return &TicketHandler{ticketStore: ts}
}

func generateSecureTicket() string {
	b := make([]byte, 32)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func (h *TicketHandler) GetWSTicket(c *gin.Context) {
	accountID, err := middleware.ExtractAccountID(c)
	if err != nil {
		apierror.DetectAndSendError(c, err)
		return
	}

	ticket := generateSecureTicket()

	// Guardamos el ticket en la store
	h.ticketStore.Store(ticket, accountID, 10*time.Second)

	c.JSON(http.StatusOK, gin.H{"ticket": ticket})
}
