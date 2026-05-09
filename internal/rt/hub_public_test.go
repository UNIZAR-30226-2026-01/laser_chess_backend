package rt

import (
	"sync"
	"testing"
	"time"

	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/api/apierror"
)

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func newPublicHub() *PublicHub {
	return NewPublicHub()
}

// newRequest construye un MatchRequest listo para usar en tests.
// Los canales son todos buffered para que NotifyMatch no bloquee.
func newRequest(accountID int64, elo, rd int32, gameMode, ranked int) *MatchRequest {
	return &MatchRequest{
		PlayerClient: stubClient(accountID),
		PlayerELO:    elo,
		PlayerRD:     rd,
		GameMode:     gameMode,
		Ranked:       ranked,
		ResponseChan: make(chan *MatchmakingFound, 1),
		ErrorChan:    make(chan error, 1),
		CancelChan:   make(chan bool, 1),
		Board:        0,
	}
}

// waitMatch espera hasta timeout a que ResponseChan reciba un match.
func waitMatch(t *testing.T, req *MatchRequest, timeout time.Duration) *MatchmakingFound {
	t.Helper()
	select {
	case found := <-req.ResponseChan:
		return found
	case <-time.After(timeout):
		t.Fatalf("timeout esperando match para el jugador %d", req.PlayerClient.AccountID)
		return nil
	}
}

// waitError espera hasta timeout a que ErrorChan reciba un error.
func waitError(t *testing.T, req *MatchRequest, timeout time.Duration) error {
	t.Helper()
	select {
	case err := <-req.ErrorChan:
		return err
	case <-time.After(timeout):
		t.Fatalf("timeout esperando error para el jugador %d", req.PlayerClient.AccountID)
		return nil
	}
}

// ---------------------------------------------------------------------------
// NewPublicHub
// ---------------------------------------------------------------------------

func TestNewPublicHub_QueuesInitialized(t *testing.T) {
	ph := newPublicHub()

	if ph.rankingOrCasualQueues[0] == nil {
		t.Error("rankingOrCasualQueues[0] (casual) es nil")
	}
	if ph.rankingOrCasualQueues[1] == nil {
		t.Error("rankingOrCasualQueues[1] (ranked) es nil")
	}
	if ph.playersInQueue == nil {
		t.Error("playersInQueue es nil")
	}
}

// ---------------------------------------------------------------------------
// EnterPlayersInQueue / ExitPlayersInQueue
// ---------------------------------------------------------------------------

func TestEnterPlayersInQueue_FirstTime_ReturnsTrue(t *testing.T) {
	ph := newPublicHub()
	if !ph.EnterPlayersInQueue(1) {
		t.Error("el primer EnterPlayersInQueue debería devolver true")
	}
}

func TestEnterPlayersInQueue_Duplicate_ReturnsFalse(t *testing.T) {
	ph := newPublicHub()
	ph.EnterPlayersInQueue(1)
	if ph.EnterPlayersInQueue(1) {
		t.Error("el segundo EnterPlayersInQueue del mismo jugador debería devolver false")
	}
}

func TestExitPlayersInQueue_AllowsReentry(t *testing.T) {
	ph := newPublicHub()
	ph.EnterPlayersInQueue(1)
	ph.ExitPlayersInQueue(1)
	if !ph.EnterPlayersInQueue(1) {
		t.Error("tras ExitPlayersInQueue el jugador debería poder volver a entrar")
	}
}

func TestExitPlayersInQueue_NonExistent_NoOp(t *testing.T) {
	ph := newPublicHub()
	ph.ExitPlayersInQueue(999) // no debe entrar en pánico
}

// ---------------------------------------------------------------------------
// AddPlayerToMatchmaking — jugador duplicado
// ---------------------------------------------------------------------------

func TestAddPlayerToMatchmaking_DuplicatePlayer_ReturnsError(t *testing.T) {
	ph := newPublicHub()

	req1 := newRequest(1, 1000, 50, 0, 0)
	req2 := newRequest(1, 1000, 50, 0, 0) // mismo AccountID

	go ph.AddPlayerToMatchmaking(req1)
	time.Sleep(20 * time.Millisecond)
	go ph.AddPlayerToMatchmaking(req2)

	err := waitError(t, req2, time.Second)
	if err != apierror.ErrAlreadyInQueue {
		t.Errorf("se esperaba ErrAlreadyInQueue, se obtuvo: %v", err)
	}

	req1.CancelChan <- true
}

// ---------------------------------------------------------------------------
// Match inmediato (mismo bracket, sin esperar ticker)
//
// Los jugadores entran de uno en uno con un sleep entre cada uno.
// Esto es necesario porque el match inmediato solo ocurre si el primer
// jugador ya ha ejecutado AddToQueue antes de que llegue el segundo.
// Si entran en paralelo ambos ven la cola vacía en CheckBracket,
// pasan a AddToQueue y quedan esperando el ticker de 3s.
// ---------------------------------------------------------------------------

func TestAddPlayerToMatchmaking_ImmediateMatch_SameBracket(t *testing.T) {
	ph := newPublicHub()

	req1 := newRequest(1, 1000, 0, 0, 0)
	req2 := newRequest(2, 1000, 0, 0, 0)

	go ph.AddPlayerToMatchmaking(req1)
	time.Sleep(20 * time.Millisecond)
	go ph.AddPlayerToMatchmaking(req2)

	found1 := waitMatch(t, req1, time.Second)
	found2 := waitMatch(t, req2, time.Second)

	if found1.Client.AccountID != 2 {
		t.Errorf("req1 debería haber emparejado con el jugador 2, obtuvo %d", found1.Client.AccountID)
	}
	if found2.Client.AccountID != 1 {
		t.Errorf("req2 debería haber emparejado con el jugador 1, obtuvo %d", found2.Client.AccountID)
	}
}

func TestAddPlayerToMatchmaking_ImmediateMatch_GameModesIsolated(t *testing.T) {
	ph := newPublicHub()

	// req1 y req2 tienen GameMode distinto → no deben matchear entre sí.
	// req3 tiene el mismo GameMode que req1 → sí deben matchear.
	//
	// NOTA: si se ejecuta con -race el detector reporta un bug real en
	// hub_public.go: CheckBracket (línea 158) lee rankingOrCasualQueues
	// sin ph.mu mientras CheckCreatedQueue escribe en ese mapa bajo ph.mu.
	// Hay que proteger esa lectura en el código productivo.
	req1 := newRequest(1, 1000, 0, 0, 0) // gameMode 0
	req2 := newRequest(2, 1000, 0, 1, 0) // gameMode 1
	req3 := newRequest(3, 1000, 0, 0, 0) // gameMode 0

	go ph.AddPlayerToMatchmaking(req1)
	time.Sleep(20 * time.Millisecond)
	go ph.AddPlayerToMatchmaking(req2)
	time.Sleep(20 * time.Millisecond)
	go ph.AddPlayerToMatchmaking(req3)

	found1 := waitMatch(t, req1, time.Second)
	found3 := waitMatch(t, req3, time.Second)

	if found1.Client.AccountID != 3 && found3.Client.AccountID != 1 {
		t.Errorf("req1 y req3 deberían haberse emparejado entre sí")
	}

	select {
	case <-req2.ResponseChan:
		t.Error("req2 (gameMode 1) no debería haber matcheado con nadie")
	default:
	}

	req2.CancelChan <- true
}

func TestAddPlayerToMatchmaking_ImmediateMatch_RankedVsCasualIsolated(t *testing.T) {
	ph := newPublicHub()

	req1 := newRequest(1, 1000, 0, 0, 0) // casual
	req2 := newRequest(2, 1000, 0, 0, 1) // ranked
	req3 := newRequest(3, 1000, 0, 0, 0) // casual → matchea con req1

	go ph.AddPlayerToMatchmaking(req1)
	time.Sleep(20 * time.Millisecond)
	go ph.AddPlayerToMatchmaking(req2)
	time.Sleep(20 * time.Millisecond)
	go ph.AddPlayerToMatchmaking(req3)

	waitMatch(t, req1, time.Second)
	waitMatch(t, req3, time.Second)

	select {
	case <-req2.ResponseChan:
		t.Error("req2 (ranked) no debería haber matcheado con jugadores casual")
	default:
	}

	req2.CancelChan <- true
}

// ---------------------------------------------------------------------------
// ELO bracket — jugadores en brackets distintos no matchean inmediatamente
// ---------------------------------------------------------------------------

func TestAddPlayerToMatchmaking_DifferentBrackets_NoImmediateMatch(t *testing.T) {
	ph := newPublicHub()

	req1 := newRequest(1, 500, 0, 0, 0)
	req2 := newRequest(2, 5000, 0, 0, 0)

	go ph.AddPlayerToMatchmaking(req1)
	go ph.AddPlayerToMatchmaking(req2)

	time.Sleep(100 * time.Millisecond)

	select {
	case <-req1.ResponseChan:
		t.Error("jugadores con ELO muy distinto no deberían matchear inmediatamente")
	default:
	}

	req1.CancelChan <- true
	req2.CancelChan <- true
}

// ---------------------------------------------------------------------------
// Cancel
// ---------------------------------------------------------------------------

func TestAddPlayerToMatchmaking_Cancel_RemovesFromQueue(t *testing.T) {
	ph := newPublicHub()

	req := newRequest(1, 1000, 0, 0, 0)
	go ph.AddPlayerToMatchmaking(req)
	time.Sleep(20 * time.Millisecond)

	req.CancelChan <- true
	time.Sleep(20 * time.Millisecond)

	// Tras cancelar el jugador debe poder volver a entrar en cola
	if !ph.EnterPlayersInQueue(1) {
		t.Error("tras cancelar el matchmaking el jugador debería poder volver a entrar en cola")
	}
	ph.ExitPlayersInQueue(1)
}

// ---------------------------------------------------------------------------
// Varios pares secuenciales
// ---------------------------------------------------------------------------

func TestAddPlayerToMatchmaking_MultipleSequentialPairs(t *testing.T) {
	ph := newPublicHub()

	// 5 pares: el primero de cada par entra y espera en cola,
	// el segundo lo encuentra inmediatamente.
	const pairs = 5
	for i := range pairs {
		p1 := newRequest(int64(i*2+1), 1000, 0, 0, 0)
		p2 := newRequest(int64(i*2+2), 1000, 0, 0, 0)

		go ph.AddPlayerToMatchmaking(p1)
		time.Sleep(10 * time.Millisecond)
		go ph.AddPlayerToMatchmaking(p2)

		waitMatch(t, p1, time.Second)
		waitMatch(t, p2, time.Second)
	}
}

// ---------------------------------------------------------------------------
// Concurrencia — race detector
// ---------------------------------------------------------------------------

func TestPublicHub_Concurrent_EnterExit_NoDataRace(t *testing.T) {
	ph := newPublicHub()

	var wg sync.WaitGroup
	for i := range 100 {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			id := int64(i)
			ph.EnterPlayersInQueue(id)
			ph.ExitPlayersInQueue(id)
		}(i)
	}
	wg.Wait()
}

func TestPublicHub_Concurrent_Matchmaking_NoDataRace(t *testing.T) {
	// Lanza 10 pares en paralelo entre sí pero con un pequeño sleep dentro
	// de cada par para garantizar match inmediato sin depender del ticker.
	// Ejecutar con: go test -race ./internal/rt/...
	ph := newPublicHub()

	const pairs = 10
	var wg sync.WaitGroup
	wg.Add(pairs * 2)

	for i := range pairs {
		p1 := newRequest(int64(i*2+1), 1000, 0, 0, 0)
		p2 := newRequest(int64(i*2+2), 1000, 0, 0, 0)

		go func() {
			defer wg.Done()
			ph.AddPlayerToMatchmaking(p1)
		}()

		time.Sleep(5 * time.Millisecond)

		go func() {
			defer wg.Done()
			ph.AddPlayerToMatchmaking(p2)
		}()

		go func() { waitMatch(t, p1, 2*time.Second) }()
		go func() { waitMatch(t, p2, 2*time.Second) }()
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("timeout global en test de concurrencia")
	}
}
