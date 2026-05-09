package rt

import (
	"context"
	"testing"
	"time"

	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/api/match"
	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/game"
	"github.com/stretchr/testify/assert"
)

// ==========================================
// MOCK DEL SERVICIO PARA AISLAR LA BDD
// ==========================================
type MockMatchService struct{}

func (m *MockMatchService) FinalizeMatch(ctx context.Context, summary match.MatchSummaryDTO) (*match.MatchRewardsDTO, error) {
	return &match.MatchRewardsDTO{
		P1XPDiff: 100, P1MoneyDiff: 10, P1EloDiff: 15,
		P2XPDiff: 50, P2MoneyDiff: 5, P2EloDiff: -15,
	}, nil
}

// ==========================================
// HELPERS
// ==========================================
func createTestClient(id int64) *Client {
	c := &Client{}
	c.AccountID = id
	// Canales generosos para que el motor escupa mensajes sin bloquearse
	c.Send = make(chan game.ResponseToRoom, 50)
	c.ToRoom = make(chan ClientSocketMessage, 50)
	c.Reconnect = make(chan bool, 50)
	c.Done = make(chan struct{})
	c.Online = true
	return c
}

// waitMsg consume mensajes del canal hasta que encuentra el tipo deseado.
// Ignora silenciosamente mensajes intermedios (como InitialState o State extra).
func waitMsg(t *testing.T, ch chan game.ResponseToRoom, expectedType game.GameMessageType) game.ResponseToRoom {
	t.Helper()
	timeout := time.After(2 * time.Second)
	for {
		select {
		case msg := <-ch:
			if msg.Type == expectedType {
				return msg
			}
		case <-timeout:
			t.Fatalf("Timeout esperando el mensaje de tipo: %s", expectedType)
			return game.ResponseToRoom{}
		}
	}
}

// ==========================================
// TEST
// ==========================================
func TestRoom_CompleteFlow(t *testing.T) {
	p1 := createTestClient(10)
	p2 := createTestClient(20)

	registry := NewMatchRegistry()
	mockDB := &MockMatchService{}

	gameInfo := &game.GameInfo{
		BoardType:     game.ACE,
		TimeBase:      30000,
		TimeIncrement: 1000,
		MatchType:     "PRIVATE",
	}

	room := &Room{}
	room.InitRoom(p1, p2, mockDB, true, gameInfo, registry)

	// 1. INICIO DE PARTIDA
	msgStartP1 := waitMsg(t, p1.Send, game.MatchStart)
	assert.Equal(t, "20", msgStartP1.Content, "P1 debería recibir el ID de P2")
	waitMsg(t, p2.Send, game.MatchStart) // Aseguramos que P2 también lo recibe

	// 2. MOVIMIENTO
	p1.ToRoom <- ClientSocketMessage{Type: string(game.Move), Content: "Tc1:c2"}

	msgMoveP2 := waitMsg(t, p2.Send, game.Move)
	assert.Contains(t, msgMoveP2.Content, "Tc1:c2")

	// 3. PAUSAS
	p1.ToRoom <- ClientSocketMessage{Type: string(game.Pause)}
	waitMsg(t, p2.Send, game.PauseRequest)

	p2.ToRoom <- ClientSocketMessage{Type: string(game.PauseReject)}
	waitMsg(t, p1.Send, game.PauseReject)

	// 4. DESCONEXIÓN
	p1.ToRoom <- ClientSocketMessage{Type: string(game.EOC)}
	waitMsg(t, p2.Send, game.Disconnection)

	// 5. RECONEXIÓN
	p1Nuevo := createTestClient(10)
	room.ReconnectChan <- ReconnectionInfo{NewClient: p1Nuevo}

	// El viejo debe recibir EOC
	waitMsg(t, p1.Send, game.EOC)
	// El rival se entera de que ha vuelto
	waitMsg(t, p2.Send, game.Reconnection)
	// El nuevo recibe la orden de reconexión con el ID rival
	msgRecon := waitMsg(t, p1Nuevo.Send, game.Reconnection)
	assert.Equal(t, "20", msgRecon.Content)

	// 6. FIN DE PARTIDA Y REWARDS
	room.Game.ToRoom <- game.ResponseToRoom{
		Type:    game.End,
		Content: "10",
		Extra:   "CHECKMATE",
	}

	// Como el End lanza primero un mensaje End y luego los Rewards, usamos el helper
	// para pescar directamente los Rewards y comprobar el Mock
	msgRewards := waitMsg(t, p1Nuevo.Send, game.Rewards)
	assert.Equal(t, "100", msgRewards.Content)
	assert.Equal(t, "10", msgRewards.Extra)
}
