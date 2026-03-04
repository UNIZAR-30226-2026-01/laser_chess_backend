package rt

// fichero que se encarga de gestionar las conexiones iniciales en
// partidas privadas
// guarda la info de que invitaciones a partidas se han hecho
// y conecta a los dos clientes cuando se acepta la partida

import "sync"

type PrivateHub struct {
	registry *MatchRegistry

	// ID del receptor -> Lista de IDs de quienes le han retado
	pendingChallenges map[int64][]int64

	// ID del usuario -> Puntero al Cliente que ya abrió el socket
	waitingPlayers map[int64]*Client

	mu sync.RWMutex
}

// Crea un hub para partidas privadas
func NewPrivateHub(r *MatchRegistry) *PrivateHub {
	return &PrivateHub{
		registry:          r,
		pendingChallenges: make(map[int64][]int64),
		waitingPlayers:    make(map[int64]*Client),
	}
}

//TODO: todo
