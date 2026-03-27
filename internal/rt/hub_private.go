package rt

// fichero que se encarga de gestionar las conexiones iniciales en
// partidas privadas
// guarda la info de que invitaciones a partidas se han hecho
// y conecta a los dos clientes cuando se acepta la partida

// Para poder acceder y consultar de forma concurrente los datos desde muchos
// puntos a la vez, utilizamos sharding.

import (
	"hash/maphash"
	"slices"
	"sync"

	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/api/apierror"
	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/game"
)

const numShards = 32

// Info inicial de la partida que se creará
type ChallengeInfo struct {
	ChallengerClient *Client
	ChallengedId     int64
	Board            game.Board_T
	StartingTime     int
	TimeIncrement    int
}

// Fragmento que contiene una porción de los mapas del hub
type challengeShard struct {
	// ID del receptor -> Lista de IDs de quienes le han retado
	pendingChallenges map[int64][]int64

	// ID del usuario -> Puntero a la info del reto que ha creado ese user
	waitingChallenges map[int64]*ChallengeInfo

	mu sync.RWMutex
}

type PrivateHub struct {
	// registro de partidas activas
	registry *MatchRegistry

	shards [numShards]*challengeShard
}

// semilla del hasher
var seed = maphash.MakeSeed()

// Devuelve el índice del shard correspondiente a un ID de usuario
func shardIndex(userID int64) int {
	h := maphash.Comparable(seed, userID)
	return int(h % numShards)
}

// Devuelve el shard correspondiente a un ID de usuario
func (ph *PrivateHub) shard(userID int64) *challengeShard {
	return ph.shards[shardIndex(userID)]
}

// Crea un hub para partidas privadas
func NewPrivateHub(r *MatchRegistry) *PrivateHub {
	ph := &PrivateHub{registry: r}
	for i := range ph.shards {
		ph.shards[i] = &challengeShard{
			pendingChallenges: make(map[int64][]int64),
			waitingChallenges: make(map[int64]*ChallengeInfo),
		}
	}
	return ph
}

// Añade al PrivateHub un nuevo reto de un user a otro
func (ph *PrivateHub) CreateChallenge(
	challenger, challenged int64,
	info *ChallengeInfo,
) error {
	indexChallenged := shardIndex(challenged)
	sChallenged := ph.shards[indexChallenged]

	indexChallenger := shardIndex(challenger)
	sChallenger := ph.shards[indexChallenger]

	// Lógica un poco complicada para hacer los locks y unlocks de forma
	// consistente (de mayor a menor) para que no haya deadlocks
	if indexChallenged == indexChallenger {
		sChallenged.mu.Lock()
		defer sChallenged.mu.Unlock()
	} else {
		first, second := indexChallenged, indexChallenger
		if first > second {
			first, second = second, first
		}
		ph.shards[first].mu.Lock()
		defer ph.shards[first].mu.Unlock()

		ph.shards[second].mu.Lock()
		defer ph.shards[second].mu.Unlock()
	}

	// Comprobar que no es un reto duplicado
	if slices.Contains(sChallenged.pendingChallenges[challenged], challenger) {
		return apierror.ErrAlreadyExists
	}
	// Añadir challenge al retado
	sChallenged.pendingChallenges[challenged] =
		append(sChallenged.pendingChallenges[challenged], challenger)

	// Guardar info del challenge
	sChallenger.waitingChallenges[challenger] = info

	return nil
}

// Devuelve la lista de IDs que han retado al usuario dado
// Devuelve nil en caso de que no haya retos
func (ph *PrivateHub) GetChallenges(challenged int64) []int64 {
	s := ph.shard(challenged)
	s.mu.RLock()
	defer s.mu.RUnlock()

	ids := s.pendingChallenges[challenged]
	if len(ids) == 0 {
		return nil
	}

	// Devolver copia
	result := make([]int64, len(ids))
	copy(result, ids)

	return result
}

// Devuelve la info del reto creado por el challenger, o nil si no existe
func (ph *PrivateHub) GetChallengeInfo(challenger int64) *ChallengeInfo {
	s := ph.shard(challenger)
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.waitingChallenges[challenger]
}

// Elimina un reto, borra la info del challenger y lo quita de la lista del challenged
func (ph *PrivateHub) RemoveChallenge(challenger, challenged int64) {
	indexChallenger := shardIndex(challenger)
	sChallenger := ph.shards[indexChallenger]

	indexChallenged := shardIndex(challenged)
	sChallenged := ph.shards[indexChallenged]

	// Hacer lock en orden para evitar deadlocks
	if indexChallenger == indexChallenged {
		sChallenger.mu.Lock()
		defer sChallenger.mu.Unlock()
	} else {
		first, second := indexChallenger, indexChallenged
		if first > second {
			first, second = second, first
		}
		ph.shards[first].mu.Lock()
		defer ph.shards[first].mu.Unlock()
		ph.shards[second].mu.Lock()
		defer ph.shards[second].mu.Unlock()
	}

	// Borrar info del challenge
	delete(sChallenger.waitingChallenges, challenger)

	// Borrar challenge pendiente
	list := sChallenged.pendingChallenges[challenged]
	for i, id := range list {
		if id == challenger {
			sChallenged.pendingChallenges[challenged] = slices.Delete(list, i, i+1)
			break
		}
	}
}

// Acepta un reto y limpia el estado
// Devuelve la ChallengeInfo del reto aceptado
func (ph *PrivateHub) AcceptChallenge(challenger, challenged int64) (*ChallengeInfo, error) {
	indexChallenger := shardIndex(challenger)
	sChallenger := ph.shards[indexChallenger]

	indexChallenged := shardIndex(challenged)
	sChallenged := ph.shards[indexChallenged]

	// Hacer el lock en orden para evitar deadlocks
	if indexChallenger == indexChallenged {
		sChallenger.mu.Lock()
		defer sChallenger.mu.Unlock()
	} else {
		first, second := indexChallenger, indexChallenged
		if first > second {
			first, second = second, first
		}
		ph.shards[first].mu.Lock()
		defer ph.shards[first].mu.Unlock()
		ph.shards[second].mu.Lock()
		defer ph.shards[second].mu.Unlock()
	}

	// Coger challenge info
	info := sChallenger.waitingChallenges[challenger]
	if info == nil {
		return nil, apierror.ErrNotFound
	}
	// Borrar la challenge info
	delete(sChallenger.waitingChallenges, challenger)

	// Borrar el challenge pendiente del challenged
	list := sChallenged.pendingChallenges[challenged]
	for i, id := range list {
		if id == challenger {
			sChallenged.pendingChallenges[challenged] = slices.Delete(list, i, i+1)
			break
		}
	}

	return info, nil
}
