package rt

import (
	"sync"
	"testing"
	"time"

	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/game"
)

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func newRegistry() *MatchRegistry {
	return NewMatchRegistry()
}

// stubRoom devuelve una Room mínima con el canal ReconnectChan inicializado,
// suficiente para los tests que no ejercitan la lógica de juego.
func stubRoom() *Room {
	return &Room{
		ReconnectChan: make(chan ReconnectionInfo, 1),
	}
}

// stubClient devuelve un Client mínimo con el AccountID dado.
func stubClient(id int64) *Client {
	return &Client{
		AccountID: id,
		Send:      make(chan game.ResponseToRoom, 4),
		Reconnect: make(chan bool, 1),
	}
}

// ---------------------------------------------------------------------------
// NewMatchRegistry
// ---------------------------------------------------------------------------

func TestNewMatchRegistry_IsEmpty(t *testing.T) {
	r := newRegistry()

	_, found := r.GetMatch(1)
	if found {
		t.Fatal("el registry recién creado no debería tener ninguna partida")
	}
}

// ---------------------------------------------------------------------------
// RegisterMatch / GetMatch
// ---------------------------------------------------------------------------

func TestRegisterMatch_BothPlayersLinked(t *testing.T) {
	r := newRegistry()
	room := stubRoom()

	r.RegisterMatch(10, 20, room)

	got1, ok1 := r.GetMatch(10)
	got2, ok2 := r.GetMatch(20)

	if !ok1 || got1 != room {
		t.Errorf("P1 no apunta a la room correcta: ok=%v room=%v", ok1, got1)
	}
	if !ok2 || got2 != room {
		t.Errorf("P2 no apunta a la room correcta: ok=%v room=%v", ok2, got2)
	}
}

func TestGetMatch_UnknownPlayer_ReturnsFalse(t *testing.T) {
	r := newRegistry()

	room, found := r.GetMatch(999)
	if found || room != nil {
		t.Fatal("GetMatch de un jugador inexistente debería devolver (nil, false)")
	}
}

func TestRegisterMatch_OverwritesPreviousRoom(t *testing.T) {
	r := newRegistry()
	room1 := stubRoom()
	room2 := stubRoom()

	r.RegisterMatch(1, 2, room1)
	r.RegisterMatch(1, 2, room2) // sobreescribe

	got, _ := r.GetMatch(1)
	if got != room2 {
		t.Fatal("RegisterMatch debería sobreescribir la room anterior para el mismo jugador")
	}
}

// ---------------------------------------------------------------------------
// RemoveMatch
// ---------------------------------------------------------------------------

func TestRemoveMatch_BothPlayersGone(t *testing.T) {
	r := newRegistry()
	room := stubRoom()

	r.RegisterMatch(10, 20, room)
	r.RemoveMatch(10, 20)

	_, ok1 := r.GetMatch(10)
	_, ok2 := r.GetMatch(20)

	if ok1 {
		t.Error("P1 debería haberse eliminado del registry")
	}
	if ok2 {
		t.Error("P2 debería haberse eliminado del registry")
	}
}

func TestRemoveMatch_NonExistent_NoOp(t *testing.T) {
	r := newRegistry()

	// No debe entrar en pánico si los IDs no existen
	r.RemoveMatch(1, 2)
}

func TestRemoveMatch_OnlyTargetPlayers(t *testing.T) {
	r := newRegistry()
	roomA := stubRoom()
	roomB := stubRoom()

	r.RegisterMatch(1, 2, roomA)
	r.RegisterMatch(3, 4, roomB)

	r.RemoveMatch(1, 2)

	_, ok3 := r.GetMatch(3)
	_, ok4 := r.GetMatch(4)

	if !ok3 || !ok4 {
		t.Error("RemoveMatch no debería afectar a partidas de otros jugadores")
	}
}

// ---------------------------------------------------------------------------
// ReconnectClient
// ---------------------------------------------------------------------------

func TestReconnectClient_PlayerWithMatch_ReturnsTrue(t *testing.T) {
	r := newRegistry()
	room := stubRoom()

	p1 := stubClient(10)
	p2 := stubClient(20)
	room.Player1 = p1
	room.Player2 = p2

	r.RegisterMatch(p1.AccountID, p2.AccountID, room)

	// Consumidor del ReconnectChan para que ReconnectPlayer no bloquee
	done := make(chan ReconnectionInfo, 1)
	go func() {
		done <- <-room.ReconnectChan
	}()

	reconnecting := stubClient(10)
	got := r.ReconnectClient(reconnecting)

	if !got {
		t.Fatal("ReconnectClient debería devolver true cuando el jugador tiene partida activa")
	}

	select {
	case info := <-done:
		if info.NewClient.AccountID != reconnecting.AccountID {
			t.Errorf("ReconnectChan recibió AccountID=%d, se esperaba %d",
				info.NewClient.AccountID, reconnecting.AccountID)
		}
	case <-time.After(time.Second):
		t.Fatal("timeout esperando señal en ReconnectChan")
	}
}

func TestReconnectClient_PlayerWithoutMatch_ReturnsFalse(t *testing.T) {
	r := newRegistry()
	player := stubClient(99)

	got := r.ReconnectClient(player)
	if got {
		t.Fatal("ReconnectClient debería devolver false cuando el jugador no tiene partida activa")
	}
}

// ---------------------------------------------------------------------------
// Concurrencia
// ---------------------------------------------------------------------------

func TestRegistry_Concurrent_NoDataRace(t *testing.T) {
	// Ejecutar con: go test -race ./internal/rt/...
	r := newRegistry()

	const goroutines = 50
	var wg sync.WaitGroup
	wg.Add(goroutines)

	for i := range goroutines {
		go func(i int) {
			defer wg.Done()

			p1 := int64(i * 2)
			p2 := int64(i*2 + 1)
			room := stubRoom()

			r.RegisterMatch(p1, p2, room)

			got, ok := r.GetMatch(p1)
			if !ok || got != room {
				// No llamamos t.Fatal desde goroutines; usamos t.Errorf
				t.Errorf("goroutine %d: GetMatch devolvió resultado inesperado", i)
			}

			r.RemoveMatch(p1, p2)

			_, stillThere := r.GetMatch(p1)
			if stillThere {
				t.Errorf("goroutine %d: el jugador sigue en el registry tras RemoveMatch", i)
			}
		}(i)
	}

	wg.Wait()
}

func TestRegistry_Concurrent_MixedOps_NoDataRace(t *testing.T) {
	// Mezcla de lecturas y escrituras simultáneas sobre las mismas claves
	r := newRegistry()
	room := stubRoom()
	r.RegisterMatch(1, 2, room)

	var wg sync.WaitGroup

	// Lectores
	for range 30 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			r.GetMatch(1)
			r.GetMatch(2)
		}()
	}

	// Escritor que registra y borra otras partidas
	for i := range 20 {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			p1, p2 := int64(100+i*2), int64(100+i*2+1)
			r.RegisterMatch(p1, p2, stubRoom())
			r.RemoveMatch(p1, p2)
		}(i)
	}

	wg.Wait()
}
