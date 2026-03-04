package rt

// match registry lleva un registro de todas las partidas activas.
// se usa para saber si un usuario ya esta jugando alguna partida o no

import "sync"

type MatchRegistry struct {
	// userID -> puntero a la Room activa
	activeMatches map[int64]*Room
	mu            sync.RWMutex
}

func NewMatchRegistry() *MatchRegistry {
	return &MatchRegistry{
		activeMatches: make(map[int64]*Room),
	}
}

// permite saber si un usuario ya tiene una partida
func (r *MatchRegistry) GetMatch(userID int64) (*Room, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	room, ok := r.activeMatches[userID]
	return room, ok
}

// vincula a ambos jugadores con una Room
func (r *MatchRegistry) RegisterMatch(p1ID, p2ID int64, room *Room) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.activeMatches[p1ID] = room
	r.activeMatches[p2ID] = room
}

// desvincula dos jugadores de una room.
// se usa al acabar una partida
func (r *MatchRegistry) RemoveMatch(p1ID, p2ID int64) {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.activeMatches, p1ID)
	delete(r.activeMatches, p2ID)
}
