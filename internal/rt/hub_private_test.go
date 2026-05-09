package rt

import (
	"errors"
	"sync"
	"testing"

	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/api/apierror"
	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/game"
)

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func newHub() *PrivateHub {
	return NewPrivateHub()
}

func stubInfo(challenger int64) *ChallengeInfo {
	return &ChallengeInfo{
		ChallengerClient: stubClient(challenger),
		ChallengedId:     0,
		Board:            game.Board_T(0),
		StartingTime:     300,
		TimeIncrement:    5,
		IsNewMatch:       true,
	}
}

// sameshardPair devuelve dos IDs distintos que caen en el mismo shard.
// Los busca en runtime porque la seed de maphash es aleatoria por proceso.
func sameshardPair() (int64, int64) {
	for a := int64(0); a < 10_000; a++ {
		for b := a + 1; b < 10_000; b++ {
			if shardIndex(a) == shardIndex(b) {
				return a, b
			}
		}
	}
	panic("no se encontró ningún par en el mismo shard (imposible con numShards=32)")
}

// diffShardPair devuelve dos IDs distintos que caen en shards distintos.
func diffShardPair() (int64, int64) {
	for a := int64(0); a < 10_000; a++ {
		for b := a + 1; b < 10_000; b++ {
			if shardIndex(a) != shardIndex(b) {
				return a, b
			}
		}
	}
	panic("no se encontró ningún par en shards distintos (imposible con numShards=32)")
}

// ---------------------------------------------------------------------------
// NewPrivateHub
// ---------------------------------------------------------------------------

func TestNewPrivateHub_AllShardsInitialized(t *testing.T) {
	ph := newHub()
	for i, s := range ph.shards {
		if s == nil {
			t.Fatalf("shard[%d] es nil", i)
		}
		if s.pendingChallenges == nil {
			t.Fatalf("shard[%d].pendingChallenges es nil", i)
		}
		if s.waitingChallenges == nil {
			t.Fatalf("shard[%d].waitingChallenges es nil", i)
		}
	}
}

// ---------------------------------------------------------------------------
// CreateChallenge — shards distintos
// ---------------------------------------------------------------------------

func TestCreateChallenge_DiffShards_OK(t *testing.T) {
	ph := newHub()
	challenger, challenged := diffShardPair()
	info := stubInfo(challenger)

	err := ph.CreateChallenge(challenger, challenged, info)
	if err != nil {
		t.Fatalf("CreateChallenge no debería devolver error: %v", err)
	}

	// El reto aparece en GetChallenges del challenged
	ids := ph.GetChallenges(challenged)
	if len(ids) != 1 || ids[0] != challenger {
		t.Errorf("GetChallenges devolvió %v, se esperaba [%d]", ids, challenger)
	}

	// La info se puede recuperar desde el challenger
	got := ph.GetChallengeInfo(challenger)
	if got != info {
		t.Errorf("GetChallengeInfo devolvió %v, se esperaba %v", got, info)
	}
}

// ---------------------------------------------------------------------------
// CreateChallenge — mismo shard (rama crítica del lock)
// ---------------------------------------------------------------------------

func TestCreateChallenge_SameShard_OK(t *testing.T) {
	ph := newHub()
	challenger, challenged := sameshardPair()
	info := stubInfo(challenger)

	err := ph.CreateChallenge(challenger, challenged, info)
	if err != nil {
		t.Fatalf("CreateChallenge (mismo shard) no debería devolver error: %v", err)
	}

	ids := ph.GetChallenges(challenged)
	if len(ids) != 1 || ids[0] != challenger {
		t.Errorf("GetChallenges devolvió %v, se esperaba [%d]", ids, challenger)
	}
}

// ---------------------------------------------------------------------------
// CreateChallenge — reto duplicado
// ---------------------------------------------------------------------------

func TestCreateChallenge_Duplicate_ReturnsError(t *testing.T) {
	ph := newHub()
	challenger, challenged := diffShardPair()
	info := stubInfo(challenger)

	_ = ph.CreateChallenge(challenger, challenged, info)
	err := ph.CreateChallenge(challenger, challenged, info)

	if !errors.Is(err, apierror.ErrAlreadyExists) {
		t.Errorf("se esperaba ErrAlreadyExists, se obtuvo: %v", err)
	}
}

func TestCreateChallenge_Duplicate_SameShard_ReturnsError(t *testing.T) {
	ph := newHub()
	challenger, challenged := sameshardPair()
	info := stubInfo(challenger)

	_ = ph.CreateChallenge(challenger, challenged, info)
	err := ph.CreateChallenge(challenger, challenged, info)

	if !errors.Is(err, apierror.ErrAlreadyExists) {
		t.Errorf("se esperaba ErrAlreadyExists, se obtuvo: %v", err)
	}
}

// ---------------------------------------------------------------------------
// CreateChallenge — varios challengers al mismo challenged
// ---------------------------------------------------------------------------

func TestCreateChallenge_MultipleChallengers(t *testing.T) {
	ph := newHub()

	// Encontrar tres IDs en shards distintos para simplificar
	var challengers []int64
	seen := map[int]bool{}
	for id := int64(0); len(challengers) < 3; id++ {
		if !seen[shardIndex(id)] {
			challengers = append(challengers, id)
			seen[shardIndex(id)] = true
		}
	}
	challenged := int64(99_999) // un ID distinto

	for _, c := range challengers {
		if err := ph.CreateChallenge(c, challenged, stubInfo(c)); err != nil {
			t.Fatalf("CreateChallenge(%d→%d) falló: %v", c, challenged, err)
		}
	}

	ids := ph.GetChallenges(challenged)
	if len(ids) != 3 {
		t.Errorf("se esperaban 3 retos pendientes, se obtuvieron %d: %v", len(ids), ids)
	}
}

// ---------------------------------------------------------------------------
// GetChallenges
// ---------------------------------------------------------------------------

func TestGetChallenges_Empty_ReturnsNil(t *testing.T) {
	ph := newHub()

	ids := ph.GetChallenges(42)
	if ids != nil {
		t.Errorf("se esperaba nil, se obtuvo %v", ids)
	}
}

func TestGetChallenges_ReturnsCopy(t *testing.T) {
	ph := newHub()
	challenger, challenged := diffShardPair()
	_ = ph.CreateChallenge(challenger, challenged, stubInfo(challenger))

	ids := ph.GetChallenges(challenged)
	ids[0] = -1 // mutar la copia

	// El estado interno no debe haberse alterado
	ids2 := ph.GetChallenges(challenged)
	if ids2[0] != challenger {
		t.Error("GetChallenges no devuelve una copia independiente")
	}
}

// ---------------------------------------------------------------------------
// GetChallengeInfo
// ---------------------------------------------------------------------------

func TestGetChallengeInfo_NotFound_ReturnsNil(t *testing.T) {
	ph := newHub()

	got := ph.GetChallengeInfo(42)
	if got != nil {
		t.Errorf("se esperaba nil, se obtuvo %v", got)
	}
}

// ---------------------------------------------------------------------------
// RemoveChallenge
// ---------------------------------------------------------------------------

func TestRemoveChallenge_DiffShards_ClearsState(t *testing.T) {
	ph := newHub()
	challenger, challenged := diffShardPair()
	_ = ph.CreateChallenge(challenger, challenged, stubInfo(challenger))

	ph.RemoveChallenge(challenger, challenged)

	if ph.GetChallengeInfo(challenger) != nil {
		t.Error("GetChallengeInfo debería ser nil tras RemoveChallenge")
	}
	if ids := ph.GetChallenges(challenged); len(ids) != 0 {
		t.Errorf("GetChallenges debería estar vacío tras RemoveChallenge, obtuvo %v", ids)
	}
}

func TestRemoveChallenge_SameShard_ClearsState(t *testing.T) {
	ph := newHub()
	challenger, challenged := sameshardPair()
	_ = ph.CreateChallenge(challenger, challenged, stubInfo(challenger))

	ph.RemoveChallenge(challenger, challenged)

	if ph.GetChallengeInfo(challenger) != nil {
		t.Error("GetChallengeInfo debería ser nil tras RemoveChallenge (mismo shard)")
	}
	if ids := ph.GetChallenges(challenged); len(ids) != 0 {
		t.Errorf("GetChallenges debería estar vacío tras RemoveChallenge (mismo shard)")
	}
}

func TestRemoveChallenge_NonExistent_NoOp(t *testing.T) {
	ph := newHub()
	// No debe entrar en pánico si el reto no existe
	ph.RemoveChallenge(1, 2)
}

func TestRemoveChallenge_OnlyRemovesTarget(t *testing.T) {
	ph := newHub()
	c1, challenged := diffShardPair()
	// Buscar un tercer ID en shard distinto a ambos
	var c2 int64 = 0
	for {
		c2++
		if c2 != c1 && c2 != challenged {
			break
		}
	}

	_ = ph.CreateChallenge(c1, challenged, stubInfo(c1))
	_ = ph.CreateChallenge(c2, challenged, stubInfo(c2))

	ph.RemoveChallenge(c1, challenged)

	ids := ph.GetChallenges(challenged)
	if len(ids) != 1 || ids[0] != c2 {
		t.Errorf("RemoveChallenge eliminó más de lo esperado: %v", ids)
	}
}

// ---------------------------------------------------------------------------
// TakeChallenge
// ---------------------------------------------------------------------------

func TestTakeChallenge_DiffShards_ReturnsInfoAndClearsState(t *testing.T) {
	ph := newHub()
	challenger, challenged := diffShardPair()
	info := stubInfo(challenger)
	_ = ph.CreateChallenge(challenger, challenged, info)

	got, err := ph.TakeChallenge(challenger, challenged)
	if err != nil {
		t.Fatalf("TakeChallenge no debería devolver error: %v", err)
	}
	if got != info {
		t.Errorf("TakeChallenge devolvió info incorrecta")
	}

	// El estado debe haber quedado limpio
	if ph.GetChallengeInfo(challenger) != nil {
		t.Error("GetChallengeInfo debería ser nil tras TakeChallenge")
	}
	if ids := ph.GetChallenges(challenged); len(ids) != 0 {
		t.Errorf("GetChallenges debería estar vacío tras TakeChallenge: %v", ids)
	}
}

func TestTakeChallenge_SameShard_ReturnsInfoAndClearsState(t *testing.T) {
	ph := newHub()
	challenger, challenged := sameshardPair()
	info := stubInfo(challenger)
	_ = ph.CreateChallenge(challenger, challenged, info)

	got, err := ph.TakeChallenge(challenger, challenged)
	if err != nil {
		t.Fatalf("TakeChallenge (mismo shard) no debería devolver error: %v", err)
	}
	if got != info {
		t.Errorf("TakeChallenge devolvió info incorrecta")
	}
}

func TestTakeChallenge_NotFound_ReturnsError(t *testing.T) {
	ph := newHub()
	challenger, challenged := diffShardPair()

	_, err := ph.TakeChallenge(challenger, challenged)
	if !errors.Is(err, apierror.ErrNotFound) {
		t.Errorf("se esperaba ErrNotFound, se obtuvo: %v", err)
	}
}

func TestTakeChallenge_IdempotentAfterFirst(t *testing.T) {
	ph := newHub()
	challenger, challenged := diffShardPair()
	_ = ph.CreateChallenge(challenger, challenged, stubInfo(challenger))

	_, _ = ph.TakeChallenge(challenger, challenged)
	_, err := ph.TakeChallenge(challenger, challenged) // segunda llamada

	if !errors.Is(err, apierror.ErrNotFound) {
		t.Errorf("segunda llamada a TakeChallenge debería devolver ErrNotFound, obtuvo: %v", err)
	}
}

// ---------------------------------------------------------------------------
// Concurrencia
// ---------------------------------------------------------------------------

func TestPrivateHub_Concurrent_NoDataRace(t *testing.T) {
	// Ejecutar con: go test -race ./internal/rt/...
	ph := newHub()

	const goroutines = 60
	var wg sync.WaitGroup
	wg.Add(goroutines)

	for i := range goroutines {
		go func(i int) {
			defer wg.Done()

			challenger := int64(i * 2)
			challenged := int64(i*2 + 1)
			info := stubInfo(challenger)

			// create → get → take
			if err := ph.CreateChallenge(challenger, challenged, info); err != nil {
				t.Errorf("goroutine %d: CreateChallenge falló: %v", i, err)
				return
			}

			ph.GetChallenges(challenged)
			ph.GetChallengeInfo(challenger)

			if _, err := ph.TakeChallenge(challenger, challenged); err != nil {
				t.Errorf("goroutine %d: TakeChallenge falló: %v", i, err)
			}
		}(i)
	}

	wg.Wait()
}

func TestPrivateHub_Concurrent_SameChallenged_NoDataRace(t *testing.T) {
	// Muchos challengers apuntando al mismo challenged simultáneamente
	ph := newHub()
	challenged := int64(99_999)

	const goroutines = 40
	var wg sync.WaitGroup
	wg.Add(goroutines)

	for i := range goroutines {
		go func(i int) {
			defer wg.Done()

			challenger := int64(i)
			info := stubInfo(challenger)

			_ = ph.CreateChallenge(challenger, challenged, info)
			ph.GetChallenges(challenged)
			ph.RemoveChallenge(challenger, challenged)
		}(i)
	}

	wg.Wait()
}
